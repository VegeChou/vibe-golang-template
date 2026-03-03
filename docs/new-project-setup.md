# Reuse Unified API Rules In A New Project

This guide explains how to make an LLM consistently follow the unified API contract in a fresh repository.

## 1) Bootstrap Files Into The New Project

From this repository, run:

```bash
bash scripts/bootstrap-api-rules.sh /path/to/your/new-project
```

This copies:

- `rules/unified-api.openapi.yaml` (machine source of truth)
- `rules/unified-api.human.md` (generated human-readable companion)
- `scripts/sync-api-docs.sh` (sync generator)
- `.github/workflows/verify-api-doc-sync.yml` (CI consistency check)
- `LLM_RULES.md` (shared constraints for all LLM tools)
- `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` (created or patched with rule references)

## 2) Tell LLM To Always Read The Rules

The bootstrap script already creates/patches tool-specific instruction files.
If you want to customize, keep this section:

```md
## API Rules
For any API/backend/frontend generation task, you MUST read:
- rules/unified-api.openapi.yaml
- rules/unified-api.human.md

Requirements:
- Follow ApiResponse/ErrorCode/pagination/i18n/security/versioning/testing rules exactly.
- Do not invent response envelope fields or error codes outside OpenAPI definitions.
- If API contract changes, update OpenAPI first, then run:
  bash scripts/sync-api-docs.sh
```

## 3) Run Initial Sync

Inside the new project:

```bash
bash scripts/sync-api-docs.sh
```

Then commit the imported files.

## 4) Day-to-Day Change Workflow

1. Edit `rules/unified-api.openapi.yaml`
2. Regenerate human doc:
   `bash scripts/sync-api-docs.sh`
3. Commit both files
4. CI verifies no drift (`verify-api-doc-sync.yml`)

## 5) What This Guarantees

- LLM has a single machine-readable contract to follow.
- Humans get a readable doc from the same source.
- CI blocks contract/document drift.
