# vibe-golang-template

A Go web backend starter template for fast vibe-coding delivery, with:
- Standard project layering (`cmd/internal/pkg/configs/scripts`)
- Example APIs (health check, create/list users)
- Unified API rules from `vibe-specs` (OpenAPI source + generated human doc)
- Shared constraints for Claude Code, Codex, and Gemini
- Scripted sync/check workflows and pre-commit guardrails

## Quick Start

```bash
make sync-ai
make sync-api
make update-specs
make fmt
make test
make run
```

After the service starts:

```bash
curl -s http://localhost:8080/healthz
curl -s -X POST http://localhost:8080/api/v1/users \
  -H 'content-type: application/json' \
  -H 'accept-language: en-US' \
  -d '{"name":"Alice","email":"alice@example.com"}'
curl -s http://localhost:8080/api/v1/users
```

## Project Structure

```text
.
├── AGENTS.md                  # Codex constraints (generated)
├── CLAUDE.md                  # Claude Code constraints (generated)
├── GEMINI.md                  # Gemini constraints (generated)
├── LLM_RULES.md               # Shared LLM API rule entry
├── ai/
│   ├── CONSTRAINTS.md         # Single source of truth for constraints
│   └── .constraints.sha256    # Generated checksum file
├── rules/
│   ├── unified-api.openapi.yaml
│   └── unified-api.human.md
├── docs/
│   ├── api-spec.md
│   ├── new-project-setup.md
│   └── architecture.md
├── scripts/
│   ├── sync-ai-constraints.sh
│   ├── verify-ai-constraints.sh
│   ├── sync-api-docs.sh
│   ├── install-githooks.sh
│   ├── bootstrap-api-rules.sh
│   └── check.sh
└── .githooks/pre-commit
```

## Unified API Rules Workflow

1. Edit `rules/unified-api.openapi.yaml`.
2. Run `make sync-api` (or `bash scripts/sync-api-docs.sh`).
3. Commit both OpenAPI and generated human doc.
4. CI verifies no drift via `.github/workflows/verify-api-doc-sync.yml`.

To pull latest rule assets from upstream `vibe-specs`:

```bash
make update-specs
```

Optional:

```bash
VIBE_SPECS_REF=main VIBE_SPECS_REPO=https://github.com/VegeChou/vibe-specs.git make update-specs
```

## Constraint Automation

1. Update `ai/CONSTRAINTS.md`.
2. Run `make sync-ai` to regenerate `AGENTS.md`, `CLAUDE.md`, and `GEMINI.md`.
3. Run `make check` to validate formatting, tests, API doc sync, and constraints.
4. Run `make hooks` to install pre-commit checks for local commits.

## Use This Template To Bootstrap Other Projects

Use this project or run:

```bash
bash scripts/bootstrap-api-rules.sh /path/to/your/new-project
```

This injects the same unified API rule assets and LLM rule references into the target repository.

## Environment Variables

See [configs/config.example.env](configs/config.example.env).
Key variables:
- `HTTP_ADDR` HTTP bind address.
- `I18N_FILE` i18n message catalog path (default `configs/i18n.json`).
