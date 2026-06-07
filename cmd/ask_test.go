package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/marioaer/sancho/internal/client"
)

func TestAskCommand(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Messages []client.Message `json:"messages"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		systemFound := false
		for _, m := range req.Messages {
			if m.Role == "system" && bytes.Contains([]byte(m.Content), []byte("TOON")) {
				systemFound = true
			}
		}
		if !systemFound {
			t.Error("expected MinLang system prompt")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"content":"analysis answer","usage":{"prompt_tokens":10,"completion_tokens":5},"finish_reason":"stop"}`))
	}))
	defer backend.Close()

	dir := t.TempDir()
	_ = os.WriteFile(dir+"/foo.go", []byte("package main"), 0644)
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir("/") }()

	askCmd := NewAskCmd(os.Stdout, os.Stderr)
	askCmd.SetArgs([]string{"-q", "what does this do", "-p", "./foo.go", "--json"})
	err := askCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
