# LLM Rules

For any API/backend/frontend task, MUST read:
1. `rules/unified-api.openapi.yaml`
2. `rules/unified-api.human.md`

Priority:
- If any conflict exists, `rules/unified-api.openapi.yaml` is the source of truth.

Hard requirements:
- Do not invent response envelope fields outside `ApiResponse` unless explicitly extended in OpenAPI.
- Do not invent error codes outside `ErrorCode`.
- Page pagination uses `page` and `size`; cursor pagination uses `cursor` and `limit`.
- Keep i18n behavior aligned with OpenAPI (`Accept-Language`, `lang`, default language).
- Generate or update global exception handling according to rules.

On API contract changes:
1. Update `rules/unified-api.openapi.yaml` first.
2. Run `bash scripts/sync-api-docs.sh`.
3. Commit both `rules/unified-api.openapi.yaml` and `rules/unified-api.human.md`.
