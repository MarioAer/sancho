package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilePrecedence(t *testing.T) {
	dir := t.TempDir()
	project := filepath.Join(dir, "project")
	_ = os.MkdirAll(project, 0755)
	_ = os.WriteFile(project+"/.sancho.json", []byte(`{"provider":"openrouter"}`), 0644)

	homeDir := filepath.Join(dir, "home", ".config", "sancho")
	_ = os.MkdirAll(homeDir, 0755)
	_ = os.WriteFile(homeDir+"/config.json", []byte(`{"provider":"anthropic"}`), 0644)
	_ = os.WriteFile(homeDir+"/config.json", []byte(`{"provider":"anthropic"}`), 0644)

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", filepath.Join(dir, "home"))
	defer os.Unsetenv("HOME")

	cfg1, _ := LoadFile(project)
	if cfg1.Provider != "openrouter" {
		t.Fatalf("expected project config first, got %s", cfg1.Provider)
	}

	os.Remove(project + "/.sancho.json")
	cfg2, _ := LoadFile(project)
	if cfg2.Provider != "anthropic" {
		t.Fatalf("expected home config fallback, got %s", cfg2.Provider)
	}
	_ = origHome
}

func TestFileJSONCStrip(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/.sancho.json", []byte(`{
        /* block comment */
        "model": "test-model",
        // line comment
        "max_tokens": 1000
    }`), 0644)

	cfg, err := LoadFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Model != "test-model" {
		t.Fatalf("expected test-model after JSONC strip, got %s", cfg.Model)
	}
}

func TestResolveCLIOverridesFile(t *testing.T) {
	fileCfg := Config{APIKey: "file-key", Model: "file/model"}
	envCfg := FromEnv()
	flags := CLIFlags{APIKey: "cli-key", Model: "cli/model", Provider: "anthropic"}
	s := Resolve(fileCfg, envCfg, flags, 0)

	if s.APIKey != "cli-key" {
		t.Fatalf("expected cli API key")
	}
	if s.Model != "cli/model" {
		t.Fatalf("expected cli model")
	}
	if s.Provider != "anthropic" {
		t.Fatalf("expected cli provider")
	}
}

func TestResolveBackwardsCompat(t *testing.T) {
	os.Setenv("WORKER_API_KEY", "legacy-key")
	os.Unsetenv("SANCHO_API_KEY")
	defer os.Unsetenv("WORKER_API_KEY")

	fileCfg := Config{}
	envCfg := FromEnv()
	flags := CLIFlags{}
	s := Resolve(fileCfg, envCfg, flags, 0)

	if s.APIKey != "legacy-key" {
		t.Fatalf("expected WORKER_API_KEY fallback, got %s", s.APIKey)
	}
}
