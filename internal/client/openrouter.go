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

type openRouterResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
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

	var orResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		return nil, fmt.Errorf("openrouter response decode failed: %w", err)
	}

	if len(orResp.Choices) > 0 {
		return &ChatResponse{
			Content: orResp.Choices[0].Message.Content,
			Usage:   orResp.Usage,
		}, nil
	}

	return &ChatResponse{}, nil
}
