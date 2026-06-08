package client

import (
	"context"
	"time"
)

type timeoutProvider struct {
	next    Provider
	timeout time.Duration
}

func WithTimeout(p Provider, timeout time.Duration) Provider {
	return &timeoutProvider{next: p, timeout: timeout}
}

func (t *timeoutProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()
	return t.next.ChatCompletion(ctx, req)
}
