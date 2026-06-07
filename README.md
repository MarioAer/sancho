# sancho

[![CI](https://github.com/marioaer/sancho/actions/workflows/ci.yml/badge.svg)](https://github.com/marioaer/sancho/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/marioaer/sancho)](https://goreportcard.com/report/github.com/marioaer/sancho)

CLI for delegating coding tasks to LLMs via OpenRouter, AWS Bedrock, OpenAI, or Anthropic.

## Install

```text
go install github.com/marioaer/sancho@latest
```

Or via Homebrew:

```text
brew install marioaer/tap/sancho
```

## Config

Search order: ./.sancho.json -> HOME/.config/sancho/config.json -> env vars -> defaults.

Env vars: SANCHO_API_KEY, SANCHO_BASE_URL, SANCHO_MODEL, SANCHO_PROVIDER

## Commands

```text
sancho ask -p '**/*.go' -q 'Where is the DB connection opened?'
sancho write -s 'Add a /health endpoint' -t cmd/server/health.go
```

## Using the Delegation Rules Template

Paste `templates/delegation-rules.template.md` into any LLM coding agent to enable `sancho` command usage. This gives the agent clear guidance on when to use each command.

The template works with:
- Claude Code / Cursor / Aider / Copilot Chat

## Examples

### Ask - Code Analysis

```text
# Find error handling patterns
sancho ask -p 'internal/**/*.go' -q 'How are errors handled in the client package?' --json

# Summarize a file before editing
sancho ask -p 'main.go' -q 'Summarize what this file does' --json

# Find all function usages
sancho ask -p '**/*.go' -q 'Where is NewProvider called?' --json
```

### Write - Code Generation

```text
# Generate a test file matching existing style
sancho write -s 'Unit tests for UserService.Validate method' -c internal/service/user.go -t internal/service/user_test.go

# Create a new handler
sancho write -s 'REST handler for user profile with GET and PUT' -t internal/handler/profile.go

# Generate documentation
sancho write -s 'Add go-doc comments to all exported functions in internal/client/provider.go' -c internal/client/provider.go -t internal/client/provider_doc.go
```

### With Config

```text
# Use specific model
sancho ask -p '*.go' -q 'What does this do?' --model 'deepseek/deepseek-chat'

# Override provider
sancho write -s 'Add logging' -t cmd/serve.go --provider anthropic

# Use with API key
SANCHO_API_KEY=sk-xxx sancho ask -p '*.go' -q 'How is auth handled?'
```
