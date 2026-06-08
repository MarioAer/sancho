package client

import (
	"testing"

	"github.com/marioaer/sancho/internal/config"
)

func TestNewProviderOpenRouter(t *testing.T) {
	s := config.Settings{Provider: "openrouter", APIKey: "sk-or", BaseURL: "https://openrouter.ai/api/v1", TimeoutSeconds: 120, Retry: config.RetryConfig{MaxAttempts: 3}}
	p := NewProvider(s)
	_, ok := p.(*timeoutProvider)
	if !ok {
		t.Fatalf("expected wrapped provider, got %T", p)
	}
}

func TestNewProviderAnthropic(t *testing.T) {
	s := config.Settings{Provider: "anthropic", APIKey: "sk-ant", TimeoutSeconds: 120, Retry: config.RetryConfig{MaxAttempts: 3}}
	p := NewProvider(s)
	if _, ok := p.(*timeoutProvider); !ok {
		t.Fatalf("expected wrapped provider, got %T", p)
	}
}

func TestNewProviderOpenAI(t *testing.T) {
	s := config.Settings{Provider: "openai", APIKey: "sk-openai", BaseURL: "https://api.openai.com/v1", TimeoutSeconds: 120, Retry: config.RetryConfig{MaxAttempts: 3}}
	p := NewProvider(s)
	if _, ok := p.(*timeoutProvider); !ok {
		t.Fatalf("expected wrapped provider, got %T", p)
	}
}

func TestNewProviderDefaultFallback(t *testing.T) {
	s := config.Settings{Provider: "unknown", APIKey: "sk", TimeoutSeconds: 120, Retry: config.RetryConfig{MaxAttempts: 3}}
	p := NewProvider(s)
	if _, ok := p.(*timeoutProvider); !ok {
		t.Fatalf("expected wrapped provider, got %T", p)
	}
}
