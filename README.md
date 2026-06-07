# sancho

[![CI](https://github.com/marioaer/sancho/actions/workflows/ci.yml/badge.svg)](https://github.com/marioaer/sancho/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/marioaer/sancho)](https://goreportcard.com/report/github.com/marioaer/sancho)

CLI for delegating coding tasks to LLMs via OpenRouter, AWS Bedrock, OpenAI, or Anthropic.

## Install

```text
go install github.com/marioaer/sancho@latest
```

## Config

Search order: ./.sancho.json -> HOME/.config/sancho/config.json -> env vars -> defaults.

Env vars: SANCHO_API_KEY, SANCHO_BASE_URL, SANCHO_MODEL, SANCHO_PROVIDER

## Commands

```text
sancho ask -p '**/*.go' -q 'Where is the DB connection opened?'
sancho write -s 'Add a /health endpoint' -t cmd/server/health.go
```
