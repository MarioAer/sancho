package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/marioaer/sancho/internal/apperr"
)

type OpenAI struct {
	APIKey  string
	BaseURL string
}

type openAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

func (o *OpenAI) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai request creation failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openai rate limited (%d)", resp.StatusCode), RetryAfter: retryAfter}
	}
	if resp.StatusCode >= 500 {
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openai server error (%d)", resp.StatusCode)}
	}

	var oaResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&oaResp); err != nil {
		return nil, fmt.Errorf("openai response decode failed: %w", err)
	}

	if len(oaResp.Choices) > 0 {
		return &ChatResponse{
			Content: oaResp.Choices[0].Message.Content,
			Usage:   oaResp.Usage,
		}, nil
	}

	return &ChatResponse{}, nil
}
