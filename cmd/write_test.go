package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWriteCommandRequiredFlags(t *testing.T) {
	writeCmd := NewWriteCmd(nil, nil)
	writeCmd.SetArgs([]string{})
	if err := writeCmd.Execute(); err == nil {
		t.Fatal("expected error when --spec and --target are missing")
	}
}

func TestWriteCommandWritesFile(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "test",
			"object":  "chat.completion",
			"created": 1,
			"model":   "fake/model",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]string{"role": "assistant", "content": "```go\nfunc hello() {}\n```"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 5, "completion_tokens": 6, "total_tokens": 11},
		})
	}))
	defer backend.Close()

	dir := t.TempDir()
	target := dir + "/out.go"
	_ = (os.Chdir)(dir)
	defer func() { _ = (os.Chdir)("/") }()

	os.Setenv("SANCHO_API_KEY", "test-key")    //nolint:errcheck
	os.Setenv("SANCHO_BASE_URL", backend.URL)  //nolint:errcheck
	os.Setenv("SANCHO_MODEL", "fake/model")    //nolint:errcheck
	os.Setenv("SANCHO_PROVIDER", "openrouter") //nolint:errcheck
	defer os.Unsetenv("SANCHO_API_KEY")        //nolint:errcheck
	defer os.Unsetenv("SANCHO_BASE_URL")       //nolint:errcheck
	defer os.Unsetenv("SANCHO_MODEL")          //nolint:errcheck
	defer os.Unsetenv("SANCHO_PROVIDER")       //nolint:errcheck

	writeCmd := NewWriteCmd(nil, nil)
	writeCmd.SetArgs([]string{"--spec", "write a hello function", "--target", target})
	if err := writeCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := (os.ReadFile)(target)
	if !strings.Contains(string(data), "func hello") {
		t.Errorf("expected generated code, got: %s", string(data))
	}
}

func TestWriteCommandStdoutContainsTargetAndChars(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "test",
			"object":  "chat.completion",
			"created": 1,
			"model":   "fake/model",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]string{"role": "assistant", "content": "```go\nfunc hello() {}\n```"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 5, "completion_tokens": 6, "total_tokens": 11},
		})
	}))
	defer backend.Close()

	dir := t.TempDir()
	target := dir + "/out.go"
	_ = (os.Chdir)(dir)
	defer func() { _ = (os.Chdir)("/") }()

	os.Setenv("SANCHO_API_KEY", "test-key")    //nolint:errcheck
	os.Setenv("SANCHO_BASE_URL", backend.URL)  //nolint:errcheck
	os.Setenv("SANCHO_MODEL", "fake/model")    //nolint:errcheck
	os.Setenv("SANCHO_PROVIDER", "openrouter") //nolint:errcheck
	defer os.Unsetenv("SANCHO_API_KEY")        //nolint:errcheck
	defer os.Unsetenv("SANCHO_BASE_URL")       //nolint:errcheck
	defer os.Unsetenv("SANCHO_MODEL")          //nolint:errcheck
	defer os.Unsetenv("SANCHO_PROVIDER")       //nolint:errcheck

	writeCmd := NewWriteCmd(nil, nil)
	writeCmd.SetArgs([]string{"--spec", "write something", "--target", target})
	if err := writeCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := (os.ReadFile)(target)
	if !strings.Contains(string(data), "func hello()") {
		t.Errorf("expected generated content, got: %s", string(data))
	}
}

func TestWriteCommandContextInjection(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "test",
			"object":  "chat.completion",
			"created": 1,
			"model":   "fake/model",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]string{"role": "assistant", "content": "actual output"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{"prompt_tokens": 5, "completion_tokens": 6, "total_tokens": 11},
		})
	}))
	defer backend.Close()

	dir := t.TempDir()
	target := dir + "/out.go"
	_ = (os.WriteFile)(dir+"/style.go", []byte("// Go style reference"), 0644)
	_ = (os.Chdir)(dir)
	defer func() { _ = (os.Chdir)("/") }()

	os.Setenv("SANCHO_API_KEY", "test-key")    //nolint:errcheck
	os.Setenv("SANCHO_BASE_URL", backend.URL)  //nolint:errcheck
	os.Setenv("SANCHO_MODEL", "fake/model")    //nolint:errcheck
	os.Setenv("SANCHO_PROVIDER", "openrouter") //nolint:errcheck
	defer os.Unsetenv("SANCHO_API_KEY")        //nolint:errcheck
	defer os.Unsetenv("SANCHO_BASE_URL")       //nolint:errcheck
	defer os.Unsetenv("SANCHO_MODEL")          //nolint:errcheck
	defer os.Unsetenv("SANCHO_PROVIDER")       //nolint:errcheck

	writeCmd := NewWriteCmd(nil, nil)
	writeCmd.SetArgs([]string{"--spec", "match style", "--target", target, "--context", "./style.go"})
	if err := writeCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := (os.ReadFile)(target)
	if !strings.Contains(string(data), "actual output") {
		t.Errorf("expected mock provider output, got: %s", string(data))
	}
}
