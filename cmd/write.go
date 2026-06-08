package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/marioaer/sancho/internal/client"
	"github.com/marioaer/sancho/internal/files"
	"github.com/spf13/cobra"
)

func NewWriteCmd(stdout, stderr *os.File) *cobra.Command {
	var spec string
	var contextFile string
	var target string

	cmd := &cobra.Command{
		Use:   "write",
		Short: "Generate code/docs from a spec",
		RunE: func(cmd *cobra.Command, args []string) error {
			if spec == "" {
				return fmt.Errorf("--spec is required")
			}
			if target == "" {
				return fmt.Errorf("--target is required")
			}

			settings := ResolveSettings(0)

			var prompt strings.Builder
			prompt.WriteString(spec)

			if contextFile != "" {
				results, err := files.ReadFiles(contextFile)
				if err != nil {
					return err
				}
				if len(results) > 0 && results[0].Error == nil {
					prompt.WriteString("\n<reference>")
					prompt.WriteString(results[0].Content)
					prompt.WriteString("</reference>")
				}
			}

			messages := []client.Message{
				{Role: "system", Content: "Generate clean, idiomatic code matching any reference style. Output ONLY the file contents."},
				{Role: "user", Content: prompt.String()},
			}

			p := client.NewProvider(settings)
			resp, err := p.ChatCompletion(cmd.Context(), client.ChatRequest{
				Model:     settings.Model,
				Messages:  messages,
				MaxTokens: settings.WriteMaxTokens,
			})
			if err != nil {
				return err
			}

			content := strings.TrimPrefix(resp.Content, "```")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)

			if err := os.WriteFile(target, []byte(content), 0644); err != nil {
				return err
			}

			fmt.Fprintf(stdout, "Wrote %s (%d chars)\n", target, len(content))    //nolint:errcheck
			fmt.Fprintf(stderr, "Tokens: %d prompt + %d completion = %d total\n", //nolint:errcheck
				resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
			return nil
		},
	}

	cmd.Flags().StringVarP(&spec, "spec", "s", "", "what to write (required)")
	cmd.Flags().StringVarP(&contextFile, "context", "c", "", "style reference file")
	cmd.Flags().StringVarP(&target, "target", "t", "", "output file path")
	return cmd
}
