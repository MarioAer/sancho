# sancho CLI Workflow Verification Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Verify CLI workflows meet expected behavior through static analysis and code review.

**Architecture:** Static code inspection validating each workflow path: ask command, write command, config resolution, and provider factory.

**Tech Stack:** Go, Cobra CLI

---

## Workflow 1: sancho ask Command

### Task 1: Verify ask input validation

**Files:**
- Read: `cmd/ask.go:22-28`

- [ ] **Step 1: Inspect input validation**

Check that `-q/--question` is required and returns error when missing.

- [ ] **Step 2: Verify error path**

Command should fail with `fmt.Errorf("-q / --question is required")` if question is empty.

---

### Task 2: Verify ask file glob expansion

**Files:**
- Read: `internal/files/read.go:16-60`

- [ ] **Step 1: Inspect ReadFiles function**

Verify `filepath.Glob` handles patterns, `**` triggers recursive walk via `filepath.Walk`.

- [ ] **Step 2: Verify error for no matches**

Returns `fmt.Errorf("no files match pattern: %s", pattern)` when matches is empty.

---

### Task 3: Verify ask TOON prompt enforcement

**Files:**
- Read: `cmd/ask.go:12-15, 43-46`

- [ ] **Step 1: Check TOON system prompt constant**

Verify `toonPrompt` constant exists with MinLang format instructions.

- [ ] **Step 2: Verify prompt composition**

User message combines file content + question, system message contains TOON instructions.

---

### Task 4: Verify ask output flow

**Files:**
- Read: `cmd/ask.go:59-65`

- [ ] **Step 1: Inspect --json flag behavior**

When `--json` is true, only content goes to stdout.

- [ ] **Step 2: Verify token cost output**

Tokens output to stderr in format: `Tokens: %d prompt + %d completion = %d total`

---

## Workflow 2: sancho write Command

### Task 5: Verify write input validation

**Files:**
- Read: `cmd/write.go:22-27`

- [ ] **Step 1: Inspect required flag validation**

Both `--spec` and `--target` required, return errors if missing.

- [ ] **Step 2: Confirm error messages**

Spec returns `--spec is required`, target returns `--target is required`.

---

### Task 6: Verify write context injection

**Files:**
- Read: `cmd/write.go:34-44`

- [ ] **Step 1: Inspect -c/--context flag handling**

Context file content wrapped in `<reference>...</reference>` tags.

- [ ] **Step 2: Verify error propagation**

ReadFiles error returned directly if context file cannot be read.

---

### Task 7: Verify write code fence stripping

**Files:**
- Read: `cmd/write.go:62-64`

- [ ] **Step 1: Inspect fence stripping logic**

`strings.TrimPrefix(resp.Content, "```")` and `strings.TrimSuffix(content, "```")`.

- [ ] **Step 2: Verify whitespace cleanup**

`strings.TrimSpace(content)` applied after fence removal.

---

### Task 8: Verify write file output

**Files:**
- Read: `cmd/write.go:66-73`

- [ ] **Step 1: Inspect WriteFile call**

Uses `os.WriteFile(target, []byte(content), 0644)`.

- [ ] **Step 2: Verify stdout output format**

`Wrote %s (%d chars)` format reported to stdout.

---

## Workflow 3: Config Resolution

### Task 9: Verify CLI flag binding

**Files:**
- Read: `cmd/root.go:24-28`

- [ ] **Step 1: Inspect persistent flag definitions**

Flags: `--api-key`, `--base-url`, `--model`, `--provider`, `--max-tokens`.

- [ ] **Step 2: Verify flag variable linkage**

Flags bound to package-level vars used in `ResolveSettings`.

---

### Task 10: Verify file config loading

**Files:**
- Read: `internal/config/file.go:48-60`

- [ ] **Step 1: Inspect config file search order**

Searches `./.sancho.json` then `~/.config/sancho/config.json`.

- [ ] **Step 2: Verify JSONC stripping**

`stripJSONC` removes `//` line comments and `/* */` block comments via regex.

---

### Task 11: Verify config precedence

**Files:**
- Read: `internal/config/resolve.go:27-91`

- [ ] **Step 1: Inspect Resolve function precedence**

For each field: CLI flags override file, file overrides env, env overrides defaults.

- [ ] **Step 2: Verify WORKER_* fallback**

`env.go:64-68` shows WORKER_API_KEY used when SANCHO_API_KEY is empty.

---

## Workflow 4: Provider Factory

### Task 12: Verify Provider interface

**Files:**
- Read: `internal/client/provider.go:33-36`

- [ ] **Step 1: Inspect Provider interface**

Two methods: `ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error)` and `SupportsModel(model string) bool`.

---

### Task 13: Verify provider factory mapping

**Files:**
- Read: `internal/client/provider.go:38-59`

- [ ] **Step 1: Inspect newProviderFunc switch cases**

Cases: "anthropic", "openai", "bedrock", "openrouter", default.

- [ ] **Step 2: Verify default falls to OpenRouter**

Default case returns OpenRouter with APIKey and BaseURL from settings.

---

### Task 14: Verify retry middleware

**Files:**
- Read: `internal/client/retry.go:128-130`

- [ ] **Step 1: Inspect exponential backoff**

Delay = min(2^attempt seconds, BackoffMaxSeconds).

- [ ] **Step 2: Verify RetryableError handling**

Uses `errors.As` check and respects `RetryAfter` duration.

---

### Task 15: Verify timeout middleware

**Files:**
- Read: `internal/client/timeout.go:109-113`, `internal/client/provider.go:65`

- [ ] **Step 1: Inspect timeout application**

`time.Duration(s.TimeoutSeconds) * time.Second` applied to context.

- [ ] **Step 2: Verify default timeout**

Default 120 seconds set in `resolve.go:34`.