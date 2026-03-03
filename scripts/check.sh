#!/usr/bin/env bash
set -euo pipefail

echo "[check] gofmt"
UNFORMATTED="$(gofmt -l . | tr -d '[:space:]')"
if [[ -n "$UNFORMATTED" ]]; then
  echo "gofmt required. run: make fmt"
  gofmt -l .
  exit 1
fi

echo "[check] go test"
go test ./...

echo "[check] api docs sync"
./scripts/sync-api-docs.sh
if ! git diff --exit-code -- rules/unified-api.human.md >/dev/null; then
  echo "rules/unified-api.human.md is out of sync. run: ./scripts/sync-api-docs.sh"
  exit 1
fi

echo "[check] constraints"
./scripts/verify-ai-constraints.sh

echo "all checks passed"
