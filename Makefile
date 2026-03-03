APP_NAME := vibe-golang-template

.PHONY: run fmt test check sync-ai sync-api update-specs hooks bootstrap-api-rules

run:
	go run ./cmd/server

fmt:
	gofmt -w .

test:
	go test ./...

sync-ai:
	./scripts/sync-ai-constraints.sh

sync-api:
	./scripts/sync-api-docs.sh

update-specs:
	./scripts/update-vibe-specs.sh

check:
	./scripts/check.sh

hooks:
	./scripts/install-githooks.sh

bootstrap-api-rules:
	@echo "Usage: ./scripts/bootstrap-api-rules.sh /path/to/target-project"
