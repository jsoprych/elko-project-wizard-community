<p align="center">
  <img src="web/src/logo.png" alt="elko-project-wizard" width="120"/>
</p>

<h1 align="center">elko-project-wizard</h1>

<p align="center">
  <strong>Stop wiring up AI agents from scratch. Generate a fully-configured, ready-to-code project in 60 seconds.</strong><br/>
  <em>Built by <a href="https://elko.ai">Elko.AI</a>'s Dark Software Factory</em>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/go-1.23+-00ADD8.svg" alt="Go Version"/></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License"/></a>
  <a href="https://elko.ai"><img src="https://img.shields.io/badge/Built_by-Elko.AI-7c6af7.svg" alt="Built by Elko.AI"/></a>
  <img src="https://img.shields.io/badge/Docker-Linux%20%7C%20macOS%20%7C%20Windows-2496ED.svg" alt="Docker"/>
</p>

---

## The Problem

Every AI-driven project starts the same painful way: you open a blank repo, write a half-baked `CLAUDE.md` from memory, forget half your coding standards, set up Docker wrong, and spend 45 minutes telling your AI agent things it should already know.

**There's a better way.**

---

## What It Does

**elko-project-wizard** is a web-based **AI project configurator**. Pick your directives — coding standards, tech stack, visibility policy, Docker setup — compose them into a profile, and hit **Generate**. You get a `.zip` that drops straight into your workflow:

```
my-project/
├── AGENTS.md          ← Full AI agent context (works with OpenCode, Cursor, Zed)
├── CLAUDE.md          ← @AGENTS.md shim — Claude Code reads this automatically
├── .gitignore         ← Stack-appropriate, zero noise
├── README.md          ← Ready to customize
├── Dockerfile         ← Multi-stage, production-ready (if Docker directive selected)
└── docker-compose.yml ← Local dev, health checks, named volumes
```

Then:

```bash
unzip my-project.zip && cd my-project
git init
claude   # or: opencode
```

**Your AI agent already knows your coding standards. Your Docker is already set up. You're writing code in under 60 seconds.**

---

## Quick Start

### Docker — runs on Linux, macOS, Windows

```bash
git clone https://github.com/jsoprych/elko-project-wizard-community.git
cd elko-project-wizard-community
docker compose up
```

Open **http://localhost:8080** and start composing.

### Build from source

```bash
git clone https://github.com/jsoprych/elko-project-wizard-community.git
cd elko-project-wizard-community
go build ./cmd/wizard && ./elko-project-wizard
```

---

## How It Works

```
  1. BROWSE              2. COMPOSE             3. GENERATE
  ────────────           ────────────           ─────────────
  Pick directives   →    Save as a Profile  →   Download .zip
  from categories        (reuse anytime)         and code.

  Visibility             Go + Minimal Deps       AGENTS.md ✓
  Tech Stack             + Greenfield            CLAUDE.md ✓
  Source Rules           + Docker                Dockerfile ✓
  Workflow               + Test-First            .gitignore ✓
  Docker                 = "my-api profile"      README.md ✓
```

---

## Built-in Directives

| Category | What's Included |
|----------|----------------|
| 🔒 **Visibility** | Private, Community/OSS, Dual Edition (private + pruned public) |
| 🛠 **Tech Stack** | Go, Node.js/TypeScript, Python — with idiomatic tooling defaults |
| 📏 **Source Rules** | Minimize Dependencies · Idiomatic Patterns · Function/File Size Limits · Greenfield Latest |
| 🧪 **Workflow** | Test-First Discipline with pre-commit hook enforcement |
| 🐳 **Docker** | Multi-stage Dockerfile + docker-compose (Linux/macOS/Windows) |

> **Custom directives?** Use the built-in AI Assist — it generates a structured prompt you can paste into Claude, ChatGPT, Grok, or DeepSeek, then import the result as a new directive. Set `ANTHROPIC_API_KEY` to have the wizard call Claude directly.

---

## Why Directives?

A **directive** is a focused, reusable markdown policy block that goes into your `AGENTS.md`. One directive = one concern. Compose as many as you need.

```json
{
  "id":       "rules-minimal-deps",
  "name":     "Minimize Dependencies",
  "category": "source-rules",
  "content":  "## Dependency Philosophy\n- Default to stdlib...\n- Every new dep requires justification..."
}
```

Mix and match. Save as a profile. Regenerate anytime. Your AI agents get consistent, project-aware context on every session.

---

## AI Assist — Two Modes

**Community (no API key needed):**
The wizard generates a perfectly-structured prompt for you to paste into any LLM — Claude, ChatGPT, Grok, DeepSeek. It tells the model exactly what directive JSON to return. Copy, paste, import. Done.

**Private (set `ANTHROPIC_API_KEY`):**
The wizard calls Claude directly and returns a ready-to-import directive JSON. Zero copy-paste.

---

## Architecture

```
cmd/wizard/          Go HTTP server — single binary, zero deps at runtime
internal/
  db/                SQLite via modernc/sqlite (pure Go — no CGO, Docker-friendly)
  directive/         Directive CRUD with builtin protection
  profile/           Saved directive compositions
  generator/         Zip scaffold builder
  ai/                Prompt sandwich + Claude API shell-out
  seed/              Loads built-in directives from directives/ at startup
  server/            REST API + static file serving
web/                 Vanilla JS SPA — elko design system (dark/light, no frameworks)
directives/          Built-in directive JSON files — add your own anytime
```

One binary. One SQLite file. No external services. Runs anywhere Docker runs.

---

## Roadmap

- [x] Directive system + profile composer
- [x] Zip scaffold generator (AGENTS.md, CLAUDE.md, .gitignore, Dockerfile)
- [x] AI Assist — prompt sandwich + Claude API
- [x] Docker deliverable (Linux / macOS / Windows)
- [x] Pre-commit hook: tests run before every commit
- [ ] CLI mode — `elko-wizard generate --profile go-private` (no server)
- [ ] More stacks — Rust, Java, Go+HTMX, Deno
- [ ] Ebook / tech doc sync generation
- [ ] agit integration — wire generated projects into workflow tracking
- [ ] Plugin system for org-specific directive libraries

---

## Contributing

Community edition is fully open. Add a directive, improve the UI, add a new tech stack.

```bash
git clone https://github.com/jsoprych/elko-project-wizard-community.git
cd elko-project-wizard-community
bash scripts/install-hooks.sh   # tests run before every commit
go test ./...
```

---

## Authors & Credits

**[Elko.AI](https://elko.ai) "TheMachine" Dark Software Factory**
John Soprych — <johnsoprych@gmail.com> — Creator & Maintainer

*elko-project-wizard: Because your AI agent deserves a proper briefing before the first line of code.*

---

**License:** MIT — see [LICENSE](./LICENSE)


---

## ☕ Support

If these tools are useful to you, consider supporting Elko.AI's open source work:

[![Support on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/elkoai)

*Built by [Elko.AI](https://elko.ai) × [DarkFabrik.AI](https://darkfabrik.ai) — open source AI dev tooling from the Dark Software Factory.* 🖤
