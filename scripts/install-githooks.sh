#!/usr/bin/env bash
set -euo pipefail

git config core.hooksPath .githooks
chmod +x .githooks/pre-commit
chmod +x scripts/sync-api-docs.sh

echo "Git hooks installed."
echo "Active hooks path: $(git config --get core.hooksPath)"
