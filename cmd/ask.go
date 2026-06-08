package cmd

import (
	"fmt"
	"os"

	"github.com/marioaer/sancho/internal/client"
	"github.com/marioaer/sancho/internal/files"
	"github.com/spf13/cobra"
)

const toonPrompt = `You are a code analysis assistant. Output ONLY TOON format:
- Use @path:line for file locations
- Use - for bullets
- No markdown, no prose, under 120 tokens.`

func NewAskCmd(stdout, stderr *os.File) *cobra.Command {
	var paths []string
	var question string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "ask",
		Short: "Delegate bulk reading to a cheap LLM",
		RunE: func(cmd *cobra.Command, args []string) error {
			if question == "" {
				return fmt.Errorf("-q / --question is required")
			}

			settings := ResolveSettings(0)
			p := client.NewProvider(settings)

			var allFiles []files.FileResult
			for _, pattern := range paths {
				r, err := files.ReadFiles(pattern)
				if err != nil {
					return err
				}
				allFiles = append(allFiles, r...)
			}

			prompt := files.FormatForPrompt(allFiles)
			messages := []client.Message{
				{Role: "system", Content: toonPrompt},
				{Role: "user", Content: prompt + "\n" + question},
			}

			req := client.ChatRequest{
				Model:     settings.Model,
				Messages:  messages,
				MaxTokens: settings.AskMaxTokens,
			}
			resp, err := p.ChatCompletion(cmd.Context(), req)
			if err != nil {
				return err
			}

			if jsonOutput {
				fmt.Fprintln(stdout, resp.Content) //nolint:errcheck
			} else {
				fmt.Fprintln(stdout, resp.Content)                                    //nolint:errcheck
				fmt.Fprintf(stderr, "Tokens: %d prompt + %d completion = %d total\n", //nolint:errcheck
					resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&paths, "paths", "p", nil, "files to ingest")
	cmd.Flags().StringVarP(&question, "question", "q", "", "extraction query")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON response instead of TOON")
	return cmd
}
