package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/marioaer/sancho/internal/apperr"
)

func parseRetryAfter(s string) time.Duration {
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return time.Duration(n) * time.Second
}

type OpenRouter struct {
	APIKey  string
	BaseURL string
}

func (o *OpenRouter) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openrouter request creation failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	httpReq.Header.Set("HTTP-Referer", "https://sancho-cli.dev")
	httpReq.Header.Set("X-Title", "sancho-cli")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openrouter request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openrouter rate limited (%d)", resp.StatusCode), RetryAfter: retryAfter}
	}
	if resp.StatusCode >= 500 {
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openrouter server error (%d)", resp.StatusCode)}
	}

	var raw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("openrouter response decode failed: %w", err)
	}

	choices, _ := raw["choices"].([]any)
	if len(choices) > 0 {
		choice, _ := choices[0].(map[string]any)
		message, _ := choice["message"].(map[string]any)
		content, _ := message["content"].(string)
		if content != "" {
			return &ChatResponse{
				Content: content,
				Usage: Usage{
					PromptTokens:     0,
					CompletionTokens: 0,
					TotalTokens:      0,
				},
				FinishReason: "stop",
			}, nil
		}
	}

	return &ChatResponse{
		Content: "",
		Usage: Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		FinishReason: "stop",
	}, nil
}

func (o *OpenRouter) SupportsModel(model string) bool {
	return true
}
