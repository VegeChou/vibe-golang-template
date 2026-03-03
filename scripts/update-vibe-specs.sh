#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VIBE_SPECS_REPO="${VIBE_SPECS_REPO:-https://github.com/VegeChou/vibe-specs.git}"
VIBE_SPECS_REF="${VIBE_SPECS_REF:-main}"

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

SOURCE_DIR="$TMP_DIR/vibe-specs"

echo "[update-specs] clone: $VIBE_SPECS_REPO ($VIBE_SPECS_REF)"
git clone --depth 1 --branch "$VIBE_SPECS_REF" "$VIBE_SPECS_REPO" "$SOURCE_DIR"
SOURCE_SHA="$(git -C "$SOURCE_DIR" rev-parse HEAD)"

echo "[update-specs] source sha: $SOURCE_SHA"

mkdir -p "$ROOT_DIR/rules" "$ROOT_DIR/scripts" "$ROOT_DIR/docs" "$ROOT_DIR/.github/workflows"

cp "$SOURCE_DIR/rules/unified-api.openapi.yaml" "$ROOT_DIR/rules/unified-api.openapi.yaml"
cp "$SOURCE_DIR/rules/unified-api.human.md" "$ROOT_DIR/rules/unified-api.human.md"
cp "$SOURCE_DIR/scripts/sync-api-docs.sh" "$ROOT_DIR/scripts/sync-api-docs.sh"
cp "$SOURCE_DIR/scripts/install-githooks.sh" "$ROOT_DIR/scripts/install-githooks.sh"
cp "$SOURCE_DIR/scripts/bootstrap-api-rules.sh" "$ROOT_DIR/scripts/bootstrap-api-rules.sh"
cp "$SOURCE_DIR/docs/api-spec.md" "$ROOT_DIR/docs/api-spec.md"
cp "$SOURCE_DIR/docs/new-project-setup.md" "$ROOT_DIR/docs/new-project-setup.md"
cp "$SOURCE_DIR/.github/workflows/verify-api-doc-sync.yml" "$ROOT_DIR/.github/workflows/verify-api-doc-sync.yml"

chmod +x "$ROOT_DIR/scripts/sync-api-docs.sh" "$ROOT_DIR/scripts/install-githooks.sh" "$ROOT_DIR/scripts/bootstrap-api-rules.sh"

# Ensure generated human doc matches copied OpenAPI in this repository context.
"$ROOT_DIR/scripts/sync-api-docs.sh"

echo "[update-specs] updated files from vibe-specs@$SOURCE_SHA"
echo "[update-specs] next: make sync-ai && make fmt test check"
