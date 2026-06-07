package config

import "os"

type Config struct {
	APIKey         string           `json:"api_key"`
	BaseURL        string           `json:"base_url"`
	Model          string           `json:"model"`
	Provider       string           `json:"provider"`
	OpenRouter     OpenRouterConfig `json:"openrouter"`
	Bedrock        BedrockConfig    `json:"bedrock"`
	OpenAI         OpenAIConfig     `json:"openai"`
	Anthropic      AnthropicConfig  `json:"anthropic"`
	Ask            AskConfig        `json:"ask"`
	Write          WriteConfig      `json:"write"`
	Retry          RetryConfig      `json:"retry"`
	TimeoutSeconds int              `json:"timeout_seconds"`
}

type OpenRouterConfig struct {
	APIKey      string
	HTTPReferer string
}

type BedrockConfig struct {
	Region        string
	ModelMappings map[string]string
}

type OpenAIConfig struct {
	APIKey       string
	Organization string
}

type AnthropicConfig struct {
	APIKey string
}

type AskConfig struct {
	MaxTokens int
}

type WriteConfig struct {
	MaxTokens int
}

type RetryConfig struct {
	MaxAttempts       int
	BackoffMaxSeconds int
	RetryStatusCodes  []int
}

func FromEnv() Config {
	cfg := Config{
		BaseURL:        "https://openrouter.ai/api/v1",
		Model:          "deepseek/deepseek-chat",
		Provider:       "openrouter",
		Ask:            AskConfig{MaxTokens: 8192},
		Write:          WriteConfig{MaxTokens: 16384},
		Retry:          RetryConfig{MaxAttempts: 3, BackoffMaxSeconds: 5, RetryStatusCodes: []int{429, 502, 503, 504}},
		TimeoutSeconds: 120,
	}

	if v := os.Getenv("SANCHO_API_KEY"); v != "" {
		cfg.APIKey = v
	} else if v := os.Getenv("WORKER_API_KEY"); v != "" {
		cfg.APIKey = v
	}

	if v := os.Getenv("SANCHO_BASE_URL"); v != "" {
		cfg.BaseURL = v
	} else if v := os.Getenv("WORKER_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}

	if v := os.Getenv("SANCHO_MODEL"); v != "" {
		cfg.Model = v
	} else if v := os.Getenv("WORKER_MODEL"); v != "" {
		cfg.Model = v
	}

	if v := os.Getenv("SANCHO_PROVIDER"); v != "" {
		cfg.Provider = v
	} else if v := os.Getenv("WORKER_PROVIDER"); v != "" {
		cfg.Provider = v
	}

	if v := os.Getenv("AWS_REGION"); v != "" {
		cfg.Bedrock.Region = v
	}

	return cfg
}
