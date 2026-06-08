package config

type CLIFlags struct {
	APIKey    string
	BaseURL   string
	Model     string
	Provider  string
	MaxTokens int
}

type Settings struct {
	APIKey         string
	BaseURL        string
	Model          string
	Provider       string
	MaxTokens      int
	AskMaxTokens   int
	WriteMaxTokens int
	TimeoutSeconds int
	OpenRouter     OpenRouterConfig
	Bedrock        BedrockConfig
	OpenAI         OpenAIConfig
	Anthropic      AnthropicConfig
	Retry          RetryConfig
}

func Resolve(fileCfg Config, envCfg Config, flags CLIFlags) Settings {
	s := Settings{
		BaseURL:        "https://openrouter.ai/api/v1",
		Model:          "deepseek/deepseek-chat",
		Provider:       "openrouter",
		AskMaxTokens:   8192,
		WriteMaxTokens: 16384,
		TimeoutSeconds: 120,
		Retry:          RetryConfig{MaxAttempts: 3, BackoffMaxSeconds: 5, RetryStatusCodes: []int{429, 502, 503, 504}},
	}

	// 1. Resolve basic settings (APIKey, BaseURL, Model, Provider)
	// Order: CLI Flag > Config File > Env Var > Default
	if flags.APIKey != "" {
		s.APIKey = flags.APIKey
	} else if fileCfg.APIKey != "" {
		s.APIKey = fileCfg.APIKey
	} else {
		s.APIKey = envCfg.APIKey
	}

	if flags.BaseURL != "" {
		s.BaseURL = flags.BaseURL
	} else if fileCfg.BaseURL != "" {
		s.BaseURL = fileCfg.BaseURL
	} else if envCfg.BaseURL != "" {
		s.BaseURL = envCfg.BaseURL
	}

	if flags.Model != "" {
		s.Model = flags.Model
	} else if fileCfg.Model != "" {
		s.Model = fileCfg.Model
	} else if envCfg.Model != "" {
		s.Model = envCfg.Model
	}

	if flags.Provider != "" {
		s.Provider = flags.Provider
	} else if fileCfg.Provider != "" {
		s.Provider = fileCfg.Provider
	} else if envCfg.Provider != "" {
		s.Provider = envCfg.Provider
	}

	// 2. Resolve Max Tokens
	// Command-specific limits from file take precedence over defaults.
	if fileCfg.Ask.MaxTokens > 0 {
		s.AskMaxTokens = fileCfg.Ask.MaxTokens
	}
	if fileCfg.Write.MaxTokens > 0 {
		s.WriteMaxTokens = fileCfg.Write.MaxTokens
	}

	// CLI flag overrides both if set.
	if flags.MaxTokens > 0 {
		s.AskMaxTokens = flags.MaxTokens
		s.WriteMaxTokens = flags.MaxTokens
	}

	// 3. Resolve other settings
	s.OpenRouter = fileCfg.OpenRouter
	s.Bedrock = fileCfg.Bedrock
	s.OpenAI = fileCfg.OpenAI
	s.Anthropic = fileCfg.Anthropic

	s.Retry = fileCfg.Retry
	if s.Retry.MaxAttempts == 0 {
		s.Retry = envCfg.Retry
	}

	s.TimeoutSeconds = fileCfg.TimeoutSeconds
	if s.TimeoutSeconds == 0 {
		s.TimeoutSeconds = envCfg.TimeoutSeconds
	}
	if s.TimeoutSeconds == 0 {
		s.TimeoutSeconds = 120
	}

	return s
}
