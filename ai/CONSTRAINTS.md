# Vibe Coding Constraints (Canonical)

## Scope
- Implement only changes directly related to the current task.
- Every change must be reversible, testable, and reviewable.
- Prefer small, incremental commits over large batch updates.

## Engineering Rules
- Go code must pass `gofmt` and `go test`.
- New business logic must include tests for at least one success path and one failure path.
- Handle errors explicitly; do not swallow errors.
- Public API changes must be reflected in docs in the same change.
- Do not introduce dependencies unrelated to the task.

## Collaboration Rules
- Output should include: changed files, core logic, validation commands, and risks.
- If requirements are unclear, provide a minimum viable implementation and state assumptions.
- Stop and warn before potentially destructive operations.

## Security Rules
- Never commit keys, tokens, or credentials.
- Never hardcode sensitive config in source code.
- Perform basic input validation on external inputs.

## Command Policy
- Standard validation commands: `make fmt test check`.
- Run pre-commit checks before pushing (`.githooks/pre-commit`).

## API Rules
For any API/backend/frontend generation task, MUST read:
- `LLM_RULES.md`
- `rules/unified-api.openapi.yaml`
- `rules/unified-api.human.md`

Requirements:
- Follow ApiResponse/ErrorCode/pagination/i18n/security/versioning/testing rules exactly.
- Do not invent response envelope fields or error codes outside OpenAPI definitions.
- If API contract changes, update OpenAPI first, then run:
  `bash scripts/sync-api-docs.sh`
- If upstream `vibe-specs` is updated and user asks to sync, run:
  `make update-specs` (or `bash scripts/update-vibe-specs.sh`)

## Agent-Specific Notes
- Claude Code: prefer small diffs and provide a concise diff summary.
- Codex: keep layering clear and change the smallest possible file set.
- Gemini: always include assumptions and validation outcomes.
