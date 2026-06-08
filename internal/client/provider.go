package client

import (
	"context"
	"time"

	"github.com/marioaer/sancho/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type ChatResponse struct {
	Content      string `json:"content"`
	Usage        Usage  `json:"usage"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Provider interface {
	ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

var newProviderFunc = func(s config.Settings) Provider {
	switch s.Provider {
	case "anthropic":
		return &Anthropic{APIKey: s.APIKey, BaseURL: "https://api.anthropic.com"}
	case "openai":
		return &OpenAI{APIKey: s.APIKey, BaseURL: s.BaseURL}
	case "bedrock":
		return &Bedrock{
			Region:   s.Bedrock.Region,
			Mappings: s.Bedrock.ModelMappings,
		}
	case "openrouter":
		return &OpenRouter{
			APIKey:  s.APIKey,
			BaseURL: s.BaseURL,
		}
	default:
		return &OpenRouter{
			APIKey:  s.APIKey,
			BaseURL: s.BaseURL,
		}
	}
}

func NewProvider(s config.Settings) Provider {
	p := newProviderFunc(s)
	p = WithRetry(p, config.RetryConfig(s.Retry))
	p = WithTimeout(p, time.Duration(s.TimeoutSeconds)*time.Second)
	return p
}
