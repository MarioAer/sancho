package config

import (
	"os"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/.sancho.json", []byte(`{"api_key":"file-key","model":"file/model","ask":{"max_tokens":100},"write":{"max_tokens":200}}`), 0644)
	os.Setenv("SANCHO_API_KEY", "env-key") //nolint:errcheck
	defer os.Unsetenv("SANCHO_API_KEY")    //nolint:errcheck

	fileCfg, _ := LoadFile(dir)
	envCfg := FromEnv()
	flags := CLIFlags{Model: "cli/model"}

	s := Resolve(fileCfg, envCfg, flags)
	if s.APIKey != "file-key" {
		t.Fatalf("expected file-key, got %s", s.APIKey)
	}
	if s.Model != "cli/model" {
		t.Fatalf("expected cli/model, got %s", s.Model)
	}
	if s.AskMaxTokens != 100 {
		t.Fatalf("expected 100 AskMaxTokens, got %d", s.AskMaxTokens)
	}
	if s.WriteMaxTokens != 200 {
		t.Fatalf("expected 200 WriteMaxTokens, got %d", s.WriteMaxTokens)
	}

	// Test CLI override for MaxTokens
	flags.MaxTokens = 500
	s = Resolve(fileCfg, envCfg, flags)
	if s.AskMaxTokens != 500 {
		t.Fatalf("expected 500 AskMaxTokens (CLI override), got %d", s.AskMaxTokens)
	}
	if s.WriteMaxTokens != 500 {
		t.Fatalf("expected 500 WriteMaxTokens (CLI override), got %d", s.WriteMaxTokens)
	}
}
