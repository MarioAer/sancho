# sancho CLI Workflow Verification Design

## Goal

Document and verify the expected behaviors of the `sancho` CLI from a user workflow perspective through static analysis.

## Architecture Summary

```
sancho/
├── cmd/
│   ├── root.go     # Cobra root command, config binding
│   ├── ask.go      # sancho ask -p <patterns> -q <question>
│   └── write.go    # sancho write -s <spec> -t <target> [-c <context>]
└── internal/
    ├── client/     # Provider interface + implementations
    ├── config/     # Config loading and resolution
    └── files/      # File reading with glob expansion
```

## Workflow 1: sancho ask

### Command Flow
```
User runs: sancho ask -p '**/*.go' -q 'Where is DB connection?' --json
         │
         ▼
NewAskCmd.Execute() → ResolveSettings() → NewProvider() → ChatCompletion()
         │                    │               │
         │                    │               └── Returns OpenRouter|Anthropic|OpenAI|Bedrock
         │                    │
         │                    └── Precedence: CLI flags → file → env → defaults
         │
         └── -q required; -p expands glob patterns; --json suppresses cost footer
```

### Verification Points
1. **Input validation** (`cmd/ask.go:26-28`): `-q/--question` required, returns error if missing
2. **File reading** (`internal/files/read.go:16-60`): `filepath.Glob` for patterns, `**` triggers walk
3. **Prompt formatting** (`internal/files/read.go:62-71`): Files wrapped in `<file path="...">` tags
4. **TOON prompt** (`cmd/ask.go:12-15`): System message enforces MinLang format
5. **Output flow** (`cmd/ask.go:62-65`): stdout = content, stderr = tokens unless `--json`

## Workflow 2: sancho write

### Command Flow
```
User runs: sancho write -s 'Add /health endpoint' -t cmd/server/health.go -c cmd/server/main.go
         │
         ▼
NewWriteCmd.Execute() → ReadFiles(context) → NewProvider() → ChatCompletion() → WriteFile()
         │                   │              │                    │
         │                   │              │                    └── Strip code fences
         │                   │              │
         │                   │              └── Uses WriteMaxTokens (default 16384)
         │                   │
         │                   └── Optional -c flag wraps content in <reference>
         │
         └── -s and -t required; validates both before API call
```

### Verification Points
1. **Input validation** (`cmd/write.go:22-27`): Both `--spec` and `--target` required
2. **Context injection** (`cmd/write.go:34-44`): Optional `-c` wraps file in `<reference>` tags
3. **System prompt** (`cmd/write.go:47`): System message instructs clean code output
4. **Fence stripping** (`cmd/write.go:62-64`): Removes ``` prefix/suffix
5. **File write** (`cmd/write.go:66-73`): WriteFile with 0644 permissions, chars count to stdout

## Workflow 3: Config Resolution

### Precedence Chain
```
CLI Flags → File Config → Env Vars → Defaults
    │           │            │         │
    │           │            │         └── provider=openrouter, model=deepseek/deepseek-chat
    │           │            │
    │           │            └── SANCHO_* preferred; WORKER_* fallback (env.go:64-86)
    │           │
    │           └── LoadFile checks: ./.sancho.json → ~/.config/sancho/config.json (file.go:49-50)
    │
    └── Global flags: --api-key, --base-url, --model, --provider, --max-tokens
```

### Verification Points
1. **CLI flags** (`cmd/root.go:24-28`): Bound to package-level vars
2. **File loading** (`internal/config/file.go:48-50`): JSONC stripped via regex
3. **Resolution** (`internal/config/resolve.go:27-91`): Each field follows precedence
4. **Provider fallback** (`internal/config/env.go:64-68`): WORKER_API_KEY if SANCHO_API_KEY missing

## Workflow 4: Provider Factory

### Provider Selection
```
Settings.Provider → newProviderFunc() → Provider instance
        │                    │
        │                    ├── "anthropic" → Anthropic{APIKey, BaseURL="https://api.anthropic.com"}
        │                    ├── "openai" → OpenAI{APIKey, BaseURL}
        │                    ├── "bedrock" → Bedrock{Region, Mappings}
        │                    └── "openrouter" (default) → OpenRouter{APIKey, BaseURL}
        │
        └── middleware applied: WithRetry() → WithTimeout()
```

### Verification Points
1. **Interface** (`internal/client/provider.go:33-36`): Provider with `ChatCompletion(ctx, req)` and `SupportsModel(model)`
2. **Factory** (`internal/client/provider.go:62-67`): Returns correct provider type, wraps with retry/timeout
3. **Retry middleware** (`internal/client/retry.go:128-130`): Exponential backoff, respects RetryableError
4. **Timeout middleware** (`internal/client/timeout.go:109-113`): Context deadline from TimeoutSeconds (default 120s)

## Verification Checklist

### ask command
- [ ] `-q` required validation
- [ ] `-p` glob pattern expansion works
- [ ] `**/*.go` recursive expansion works
- [ ] TOON system prompt is sent
- [ ] Response content to stdout
- [ ] Token usage to stderr (unless --json)

### write command
- [ ] `--spec` required validation
- [ ] `--target` required validation
- [ ] Context file injection works
- [ ] Code fence stripping works
- [ ] File written to target path
- [ ] Char count reported

### config
- [ ] CLI flag overrides file config
- [ ] File config overrides env vars
- [ ] Env var fallback to WORKER_* works
- [ ] Default values applied when nothing set

### providers
- [ ] OpenRouter is default provider
- [ ] Provider factory returns correct type
- [ ] Retry middleware wraps calls
- [ ] Timeout middleware applies deadline