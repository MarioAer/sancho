# sancho

[![CI](https://github.com/marioaer/sancho/actions/workflows/ci.yml/badge.svg)](https://github.com/marioaer/sancho/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/marioaer/sancho)](https://goreportcard.com/report/github.com/marioaer/sancho)

CLI for delegating coding tasks to LLMs via OpenRouter, AWS Bedrock, OpenAI, or Anthropic.

## Install

```bash
go install github.com/marioaer/sancho@latest
```

## Quick start

```bash
export SANCHO_API_KEY=sk-xxx
sancho ask -p '**/*.go' -q 'Where is the DB connection opened?'
sancho write -s 'Add a /health endpoint' -t cmd/server/health.go
```

## commands

### `sancho ask`

Read files and ask a question. Outputs TOON format by default.

```bash
sancho ask -p 'internal/**/*.go' -q 'How are errors handled?'
sancho ask -p 'main.go' -q 'Summarize what this file does'
```

Use `--json` to get raw JSON output:

```bash
sancho ask -p '**/*.go' -q 'Where is NewProvider called?' --json
```

### `sancho write`

Generate code/docs from a spec.

```bash
sancho write -s 'Unit tests for UserService.Validate' -c internal/service/user.go -t internal/service/user_test.go
sancho write -s 'REST handler for user profile with GET and PUT' -t internal/handler/profile.go
```

## Flags

| flag | shorthand | description | default |
|---|---|---|---|
| `--paths` | `-p` | files or globs to read | |
| `--question` | `-q` | extraction query | |
| `--spec` | `-s` | what to generate (required for `write`) | |
| `--target` | `-t` | output file path (required for `write`) | |
| `--context` | `-c` | style reference file | |
| `--json` | | emit JSON (suppress TOON/system prompt in `ask`) | TOON for `ask` |
| `--model` | | model override | `deepseek/deepseek-chat` |
| `--provider` | | provider override | `openrouter` |
| `--base-url` | | provider URL override | |
| `--api-key` | | API key override | |
| `--max-tokens` | | max tokens override | `8192` / `16384` |

## Config

Search order: `./.sancho.json` -> `~/.config/sancho/config.json` -> env vars -> defaults.

Env vars: `SANCHO_API_KEY`, `SANCHO_BASE_URL`, `SANCHO_MODEL`, `SANCHO_PROVIDER`
