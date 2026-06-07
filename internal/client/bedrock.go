package client

import (
	"context"
)

type Bedrock struct {
	Region   string
	Mappings map[string]string
}

func (b *Bedrock) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	return nil, nil
}

func (b *Bedrock) SupportsModel(model string) bool {
	return true
}

func (b *Bedrock) mapModel(model string) string {
	if b.Mappings == nil {
		return model
	}
	if mapped, ok := b.Mappings[model]; ok {
		return mapped
	}
	return model
}
