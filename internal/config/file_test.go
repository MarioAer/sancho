package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileJSONC(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".sancho.json")
	_ = os.WriteFile(path, []byte(`{
		// comment
		"api_key": "sk-file",
		"model": "file/model"
	}`), 0644)

	cfg, err := LoadFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "sk-file" {
		t.Fatalf("expected sk-file, got %s", cfg.APIKey)
	}
	if cfg.Model != "file/model" {
		t.Fatalf("expected file/model, got %s", cfg.Model)
	}
}

func TestLoadFileMissing(t *testing.T) {
	cfg, err := LoadFile("/nonexistent/path")
	if err != nil {
		t.Fatalf("unexpected error on missing file: %v", err)
	}
	if cfg.APIKey != "" {
		t.Fatalf("expected empty APIKey, got %s", cfg.APIKey)
	}
}
