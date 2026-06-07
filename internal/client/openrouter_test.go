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

func TestOpenRouterChatCompletion(t *testing.T) {
	var received struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk-or-test" {
			t.Errorf("expected Bearer sk-or-test, got %s", r.Header.Get("Authorization"))
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "test",
			"object":  "chat.completion",
			"created": 1,
			"model":   "test-model",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]string{"role": "assistant", "content": "answer"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15},
		})
	}))
	defer backend.Close()

	p := &OpenRouter{
		APIKey:  "sk-or-test",
		BaseURL: backend.URL,
	}
	req := ChatRequest{
		Model:     "test-model",
		Messages:  []Message{{Role: "user", Content: "hi"}},
		MaxTokens: 512,
	}
	resp, err := p.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "answer" {
		t.Fatalf("expected answer, got %s", resp.Content)
	}
	if received.Model != "test-model" {
		t.Fatalf("expected model test-model in request")
	}
}

func TestOpenRouterRateLimit(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer backend.Close()

	p := &OpenRouter{APIKey: "test", BaseURL: backend.URL}
	_, err := p.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected rate limit error")
	}
	var retryErr *apperr.RetryableError
	if !apperr.AsRetryable(err, &retryErr) {
		t.Fatalf("expected RetryableError, got %T", err)
	}
	if retryErr.RetryAfter != time.Second {
		t.Fatalf("expected 1s retry-after, got %v", retryErr.RetryAfter)
	}
}
