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
		return nil, &apperr.RetryableError{Msg: "rate limited", RetryAfter: 100 * time.Millisecond}
	}
	return &ChatResponse{Content: "ok"}, nil
}

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

func TestRetryRespectsRetryAfter(t *testing.T) {
	p := &staticProvider{err: &apperr.RetryableError{Msg: "slow down", RetryAfter: 200 * time.Millisecond}}
	wrapped := WithRetry(p, config.RetryConfig{MaxAttempts: 2, BackoffMaxSeconds: 5})

	_, err := wrapped.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected error after max attempts")
	}
}

func TestRetryRespectsLargeRetryAfter(t *testing.T) {
	p := &staticProvider{err: &apperr.RetryableError{Msg: "slow down", RetryAfter: 500 * time.Millisecond}}
	wrapped := WithRetry(p, config.RetryConfig{MaxAttempts: 2, BackoffMaxSeconds: 1})

	start := time.Now()
	_, _ = wrapped.ChatCompletion(context.Background(), ChatRequest{})
	duration := time.Since(start)

	if duration < 500*time.Millisecond {
		t.Errorf("expected wait of at least 500ms, got %v", duration)
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
