package client

import (
	"testing"

	"github.com/marioaer/sancho/internal/config"
)

func TestNewProviderOpenRouter(t *testing.T) {
	s := config.Settings{Provider: "openrouter", APIKey: "sk-or", BaseURL: "https://openrouter.ai/api/v1"}
	p := NewProvider(s)
	if _, ok := p.(*OpenRouter); !ok {
		t.Fatalf("expected *OpenRouter, got %T", p)
	}
}

func TestNewProviderAnthropic(t *testing.T) {
	s := config.Settings{Provider: "anthropic", APIKey: "sk-ant"}
	p := NewProvider(s)
	if _, ok := p.(*Anthropic); !ok {
		t.Fatalf("expected *Anthropic, got %T", p)
	}
}

func TestNewProviderOpenAI(t *testing.T) {
	s := config.Settings{Provider: "openai", APIKey: "sk-openai", BaseURL: "https://api.openai.com/v1"}
	p := NewProvider(s)
	if _, ok := p.(*OpenAI); !ok {
		t.Fatalf("expected *OpenAI, got %T", p)
	}
}

func TestNewProviderDefaultFallback(t *testing.T) {
	s := config.Settings{Provider: "unknown", APIKey: "sk"}
	p := NewProvider(s)
	if _, ok := p.(*OpenRouter); !ok {
		t.Fatalf("expected default *OpenRouter, got %T", p)
	}
}
