# Architecture Notes

- `cmd/server`: process entrypoint.
- `internal/app`: dependency wiring and route registration.
- `internal/controller`: HTTP layer, protocol translation only.
- `internal/service`: business rules.
- `internal/repository`: data persistence implementations.
- `pkg`: reusable cross-module utilities.
- `rules`: unified API contract source (`unified-api.openapi.yaml`) and generated human doc.

This template uses an in-memory repository for demonstration. Replace implementations under `internal/repository` with MySQL/PostgreSQL/Redis adapters for production use.
HTTP controllers return the unified `ApiResponse` envelope (success/code/message/lang/data/traceId/timestamp) with language resolved by `Accept-Language` then `lang` query fallback.
Request parsing follows unified conventions (`page/size` and `cursor/limit`, with default and max validation), and global panic recovery middleware maps unexpected failures to unified `ApiErrorResponse`.
Localization messages are resolved by key from `configs/i18n.json` (`I18N_FILE`), instead of hardcoding bilingual strings in handlers.
