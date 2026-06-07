package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marioaer/sancho/internal/apperr"
)

type Anthropic struct {
	APIKey  string
	BaseURL string
}

func (a *Anthropic) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":      req.Model,
		"max_tokens": req.MaxTokens,
		"messages":   req.Messages,
	})
	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.BaseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("anthropic request creation failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("anthropic rate limited (%d)", resp.StatusCode)}
	}
	if resp.StatusCode >= 500 {
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("anthropic server error (%d)", resp.StatusCode)}
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("anthropic response decode failed: %w", err)
	}

	content, _ := raw["content"].([]interface{})
	text := ""
	if len(content) > 0 {
		if block, ok := content[0].(map[string]interface{}); ok {
			text, _ = block["text"].(string)
		}
	}

	usage, _ := raw["usage"].(map[string]interface{})
	inputTok, _ := usage["input_tokens"].(float64)
	outputTok, _ := usage["output_tokens"].(float64)

	return &ChatResponse{
		Content: text,
		Usage: Usage{
			PromptTokens:     int(inputTok),
			CompletionTokens: int(outputTok),
			TotalTokens:      int(inputTok + outputTok),
		},
		FinishReason: "stop",
	}, nil
}

func (a *Anthropic) SupportsModel(model string) bool {
	return true
}
