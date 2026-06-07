package client

import (
	"testing"
)

func TestBedrockModelMapping(t *testing.T) {
	p := &Bedrock{
		Region: "us-east-1",
		Mappings: map[string]string{
			"claude-3-opus": "anthropic.claude-3-opus-20240307-v1:0",
		},
	}

	mapped := p.mapModel("claude-3-opus")
	if mapped != "anthropic.claude-3-opus-20240307-v1:0" {
		t.Fatalf("expected bedrock model id, got %s", mapped)
	}
	unmapped := p.mapModel("gpt-4")
	if unmapped != "gpt-4" {
		t.Fatalf("passthrough expected, got %s", unmapped)
	}
}

func TestBedrockRegionDefault(t *testing.T) {
	p := &Bedrock{}
	if p.Region != "" {
		t.Fatalf("expected empty default region, got %s", p.Region)
	}
}
