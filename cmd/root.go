package cmd

import (
	"os"

	"github.com/marioaer/sancho/internal/config"
	"github.com/spf13/cobra"
)

var (
	apiKey    string
	baseURL   string
	model     string
	provider  string
	maxTokens int
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "sancho",
		Short: "Delegate coding tasks to LLMs",
	}

	root.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key override")
	root.PersistentFlags().StringVar(&baseURL, "base-url", "", "Provider endpoint override")
	root.PersistentFlags().StringVar(&model, "model", "", "Model override")
	root.PersistentFlags().StringVar(&provider, "provider", "", "Provider override (openrouter, bedrock, openai, anthropic)")
	root.PersistentFlags().IntVar(&maxTokens, "max-tokens", 0, "Max tokens override")

	root.AddCommand(NewAskCmd(os.Stdout, os.Stderr))
	root.AddCommand(NewWriteCmd(os.Stdout, os.Stderr))

	return root
}

func ResolveSettings(cmdMaxTokens int) config.Settings {
	cwd, _ := os.Getwd()
	fileCfg, _ := config.LoadFile(cwd)
	envCfg := config.FromEnv()

	return config.Resolve(fileCfg, envCfg, config.CLIFlags{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		Provider:  provider,
		MaxTokens: maxTokens,
	}, cmdMaxTokens)
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
