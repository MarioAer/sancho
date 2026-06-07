package client

import (
	"context"
	"math"
	"time"

	"github.com/marioaer/sancho/internal/apperr"
	"github.com/marioaer/sancho/internal/config"
)

type retryProvider struct {
	next        Provider
	maxAttempts int
	backoffMax  time.Duration
}

func WithRetry(p Provider, retry config.RetryConfig) Provider {
	return &retryProvider{
		next:        p,
		maxAttempts: retry.MaxAttempts,
		backoffMax:  time.Duration(retry.BackoffMaxSeconds) * time.Second,
	}
}

func (r *retryProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var lastErr error
	for attempt := 0; attempt < r.maxAttempts; attempt++ {
		resp, err := r.next.ChatCompletion(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		var retryErr *apperr.RetryableError
		if !apperr.AsRetryable(lastErr, &retryErr) {
			return nil, lastErr
		}
		delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		if delay > r.backoffMax {
			delay = r.backoffMax
		}
		if retryErr.RetryAfter > 0 && retryErr.RetryAfter < delay {
			delay = retryErr.RetryAfter
		}
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, lastErr
}

func (r *retryProvider) SupportsModel(model string) bool {
	return r.next.SupportsModel(model)
}
