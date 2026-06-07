package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutMiddleware(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	p := &OpenRouter{APIKey: "test", BaseURL: backend.URL}
	wrapped := WithTimeout(p, 50*time.Millisecond)

	_, err := wrapped.ChatCompletion(context.Background(), ChatRequest{})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
