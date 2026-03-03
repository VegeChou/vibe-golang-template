# Architecture Notes

- `cmd/server`: process entrypoint.
- `internal/app`: dependency wiring and route registration.
- `internal/handler`: HTTP layer, protocol translation only.
- `internal/service`: business rules.
- `internal/repository`: data persistence implementations.
- `pkg`: reusable cross-module utilities.

This template uses an in-memory repository for demonstration. Replace implementations under `internal/repository` with MySQL/PostgreSQL/Redis adapters for production use.
