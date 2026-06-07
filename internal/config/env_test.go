package config

import (
	"os"
	"testing"
)

func TestEnvConfigReadsSANCHO(t *testing.T) {
	os.Setenv("SANCHO_API_KEY", "sk-or-test")
	os.Setenv("SANCHO_MODEL", "custom/model")
	defer os.Unsetenv("SANCHO_API_KEY")
	defer os.Unsetenv("SANCHO_MODEL")

	cfg := FromEnv()
	if cfg.APIKey != "sk-or-test" {
		t.Fatalf("expected sk-or-test, got %s", cfg.APIKey)
	}
	if cfg.Model != "custom/model" {
		t.Fatalf("expected custom/model, got %s", cfg.Model)
	}
}

func TestEnvConfigWorkerFallback(t *testing.T) {
	os.Setenv("WORKER_API_KEY", "sk-legacy")
	os.Unsetenv("SANCHO_API_KEY")
	defer os.Unsetenv("WORKER_API_KEY")

	cfg := FromEnv()
	if cfg.APIKey != "sk-legacy" {
		t.Fatalf("expected WORKER_API_KEY fallback, got %s", cfg.APIKey)
	}
}

func TestEnvConfigDefaults(t *testing.T) {
	cfg := FromEnv()
	if cfg.Provider != "openrouter" {
		t.Fatalf("expected default provider openrouter, got %s", cfg.Provider)
	}
	if cfg.Model != "deepseek/deepseek-chat" {
		t.Fatalf("expected default model")
	}
	if cfg.BaseURL != "https://openrouter.ai/api/v1" {
		t.Fatalf("expected default base URL")
	}
}
