# vibe-golang-template

A Go web backend starter template for fast vibe-coding delivery, with:
- Standard project layering (`cmd/internal/pkg/configs/scripts`)
- Example APIs (health check, create/list users)
- Shared constraints for Claude Code, Codex, and Gemini
- Scripted constraint sync and checksum validation
- Pre-commit automation for local guardrails

## Quick Start

```bash
make sync-ai
make fmt
make test
make run
```

After the service starts:

```bash
curl -s http://localhost:8080/healthz
curl -s -X POST http://localhost:8080/api/v1/users \
  -H 'content-type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
curl -s http://localhost:8080/api/v1/users
```

## Project Structure

```text
.
├── AGENTS.md                  # Codex constraints (generated)
├── CLAUDE.md                  # Claude Code constraints (generated)
├── GEMINI.md                  # Gemini constraints (generated)
├── ai/
│   ├── CONSTRAINTS.md         # Single source of truth for constraints
│   └── .constraints.sha256    # Generated checksum file
├── cmd/server/main.go
├── internal/
│   ├── app/                   # Wiring and bootstrapping
│   ├── config/                # Configuration loading
│   ├── handler/               # HTTP transport layer
│   ├── model/                 # Domain models
│   ├── repository/            # Data access layer
│   └── service/               # Business logic layer
├── pkg/response/              # Reusable response helpers
├── scripts/
│   ├── sync-ai-constraints.sh
│   ├── verify-ai-constraints.sh
│   ├── check.sh
│   └── bootstrap-hooks.sh
└── .githooks/pre-commit
```

## Constraint Automation

1. Update `ai/CONSTRAINTS.md`.
2. Run `make sync-ai` to regenerate `AGENTS.md`, `CLAUDE.md`, and `GEMINI.md`.
3. Run `make check` to validate formatting, tests, and constraint checksums.
4. Run `make hooks` to install pre-commit checks for local commits.

## Environment Variables

See [configs/config.example.env](configs/config.example.env).
