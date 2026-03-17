# AGENTS.md — elko-project-wizard (Community Edition)

AI agent directives for Claude Code, OpenCode, Cursor, and Zed.

## Project Overview

**elko-project-wizard** is a Go web app that composes modular AI agent directives
into `AGENTS.md` / `CLAUDE.md` files and generates ready-to-use project scaffold zips.

## Build & Test

```bash
go build ./...
go test ./...
bash scripts/install-hooks.sh  # installs pre-commit test hook
```

## Run

```bash
# Docker (Linux/macOS/Windows)
docker compose up

# From source
./elko-project-wizard --port 8080
```

## Architecture

```
cmd/wizard/       HTTP server entrypoint
internal/db/      SQLite (pure Go, no CGO)
internal/directive Directive CRUD
internal/profile/ Profile management
internal/generator Zip scaffold builder
internal/ai/      Prompt sandwich + Claude API assist
internal/seed/    Built-in directive loader
internal/server/  HTTP routes + handlers
web/              Vanilla JS SPA
directives/       Built-in directive JSON files
```

## Code Standards

- Max function: 80 lines. Max file: 300 lines.
- Prefer stdlib over external packages.
- Every new file gets a `_test.go` counterpart.
- Tests enforced via pre-commit hook.
- Pure Go only (no CGO).
