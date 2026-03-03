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

echo "[check] constraints"
./scripts/verify-ai-constraints.sh

echo "all checks passed"
