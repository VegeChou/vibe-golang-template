#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECKSUM_FILE="$ROOT_DIR/ai/.constraints.sha256"

if [[ ! -f "$CHECKSUM_FILE" ]]; then
  echo "missing checksum file, run: make sync-ai" >&2
  exit 1
fi

if command -v shasum >/dev/null 2>&1; then
  shasum -a 256 -c "$CHECKSUM_FILE"
elif command -v sha256sum >/dev/null 2>&1; then
  sha256sum -c "$CHECKSUM_FILE"
else
  echo "no checksum command found" >&2
  exit 1
fi

echo "constraints checksum verified"
