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
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openai rate limited (%d)", resp.StatusCode), RetryAfter: retryAfter}
	}
	if resp.StatusCode >= 500 {
		return nil, &apperr.RetryableError{Msg: fmt.Sprintf("openai server error (%d)", resp.StatusCode)}
	}

	var cr ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, fmt.Errorf("openai response decode failed: %w", err)
	}
	return &cr, nil
}

func (o *OpenAI) SupportsModel(model string) bool {
	return true
}
