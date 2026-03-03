APP_NAME := vibe-golang-template

.PHONY: run fmt test check sync-ai hooks

run:
	go run ./cmd/server

fmt:
	gofmt -w .

test:
	go test ./...

sync-ai:
	./scripts/sync-ai-constraints.sh

check:
	./scripts/check.sh

hooks:
	./scripts/bootstrap-hooks.sh
