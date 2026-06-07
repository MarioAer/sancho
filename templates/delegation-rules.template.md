# sancho delegation rules

Paste these rules into any LLM coding agent (Claude Code, Cursor, Aider, Copilot Chat) to enable `sancho` invocations.

---

## When to use `sancho ask`

Use `sancho ask` when you need to:
- Scan multiple files for a pattern or answer a specific question
- Summarize existing code before editing
- Find all usages of a function/type/constant

Example: "Where do we handle auth errors?"

```text
sancho ask -p "**/*.go" -q "Where are auth errors logged?" --json
```

---

## When to use `sancho write`

Use `sancho write` when you need to:
- Generate a new file from a specification
- Create a test file matching existing style
- Scaffold a component or package

Example: "Write a test for the user service following existing patterns"

```text
sancho write -s "Unit tests for UserService covering CRUD" -c internal/service/user.go -t internal/service/user_test.go
```

---

## Rules

1. Always pass `--json` when piping output to another tool.
2. Use `-c / --context` with `sancho write` to match project style.
3. Use glob patterns with `-p` (supports `**` for recursive).
4. Model overrides: `--model`. Provider overrides: `--provider`.
