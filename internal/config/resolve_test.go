package config

import (
	"os"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/.sancho.json", []byte(`{"api_key":"file-key","model":"file/model"}`), 0644)
	os.Setenv("SANCHO_API_KEY", "env-key")
	defer os.Unsetenv("SANCHO_API_KEY")

	fileCfg, _ := LoadFile(dir)
	envCfg := FromEnv()
	flags := CLIFlags{Model: "cli/model"}

	s := Resolve(fileCfg, envCfg, flags, 0)
	if s.APIKey != "file-key" {
		t.Fatalf("expected file-key, got %s", s.APIKey)
	}
	if s.Model != "cli/model" {
		t.Fatalf("expected cli/model, got %s", s.Model)
	}
}
