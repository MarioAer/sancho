package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marioaer/sancho/internal/apperr"
)

func TestAnthropicChatCompletion(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("expected /v1/messages, got %s", r.URL.Path)
		}
		if r.Header.Get("anthropic-version") == "" {
			t.Error("expected anthropic-version header")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"content":[{"text":"anthropic-answer","type":"text"}],"usage":{"input_tokens":5,"output_tokens":10}}`))
	}))
	defer backend.Close()

	p := &Anthropic{
		APIKey:  "sk-ant-test",
		BaseURL: backend.URL,
	}
	resp, err := p.ChatCompletion(context.Background(), ChatRequest{
		Model:    "claude-3-opus",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "anthropic-answer" {
		t.Fatalf("expected anthropic-answer, got %s", resp.Content)
	}
}

func TestAnthropicRateLimit(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "3")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"rate limit"}}`))
	}))
	defer backend.Close()

	p := &Anthropic{APIKey: "test", BaseURL: backend.URL}
	_, err := p.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	var retryErr *apperr.RetryableError
	if !apperr.AsRetryable(err, &retryErr) {
		t.Fatalf("expected RetryableError, got %T", err)
	}
	if retryErr.Msg == "" {
		t.Fatalf("expected non-empty retry error message, got %q", retryErr.Msg)
	}
}
