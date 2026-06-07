package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marioaer/sancho/internal/apperr"
)

func TestOpenAIChatCompletion(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected /chat/completions, got %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&struct{}{})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ChatResponse{
			Content:      "openai-answer",
			Usage:        Usage{TotalTokens: 20},
			FinishReason: "stop",
		})
	}))
	defer backend.Close()

	p := &OpenAI{
		APIKey:  "sk-openai",
		BaseURL: backend.URL,
	}
	resp, err := p.ChatCompletion(context.Background(), ChatRequest{
		Model:    "gpt-4",
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "openai-answer" {
		t.Fatalf("expected openai-answer, got %s", resp.Content)
	}
}

func TestOpenAIRateLimit(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "2")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer backend.Close()

	p := &OpenAI{APIKey: "test", BaseURL: backend.URL}
	_, err := p.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	var retryErr *apperr.RetryableError
	if !apperr.AsRetryable(err, &retryErr) {
		t.Fatalf("expected RetryableError, got %T", err)
	}
	if retryErr.RetryAfter != 2*time.Second {
		t.Fatalf("expected 2s retry-after, got %v", retryErr.RetryAfter)
	}
}
