package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/marioaer/sancho/internal/apperr"
	"github.com/marioaer/sancho/internal/config"
)

type failingProvider struct {
	attempts int
	max      int
}

func (f *failingProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	f.attempts++
	if f.attempts < f.max {
		return nil, &apperr.RetryableError{Msg: "rate limited", RetryAfter: time.Second}
	}
	return &ChatResponse{Content: "ok"}, nil
}

func (f *failingProvider) SupportsModel(m string) bool { return true }

func TestRetryMiddleware(t *testing.T) {
	p := &failingProvider{max: 3}
	wrapped := WithRetry(p, config.RetryConfig{MaxAttempts: 3, BackoffMaxSeconds: 5, RetryStatusCodes: []int{429, 502, 503, 504}})

	_, err := wrapped.ChatCompletion(context.Background(), ChatRequest{})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if p.attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", p.attempts)
	}
}

func TestRetryNoRetryOnClientError(t *testing.T) {
	p := &staticProvider{err: fmt.Errorf("bad request")}
	wrapped := WithRetry(p, config.RetryConfig{MaxAttempts: 3, BackoffMaxSeconds: 5, RetryStatusCodes: []int{429, 502, 503, 504}})

	_, err := wrapped.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected error on non-retryable failure")
	}
	if p.calls != 1 {
		t.Fatalf("expected no retry on client error, got %d calls", p.calls)
	}
}

type staticProvider struct {
	err   error
	calls int
}

func (s *staticProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	s.calls++
	return nil, s.err
}

func (s *staticProvider) SupportsModel(m string) bool { return true }
