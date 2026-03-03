#!/usr/bin/env bash
set -euo pipefail

# Bootstrap unified API rule assets into a target project.
# Usage:
#   bash scripts/bootstrap-api-rules.sh /path/to/target-project

if [[ $# -ne 1 ]]; then
  echo "Usage: bash scripts/bootstrap-api-rules.sh /path/to/target-project" >&2
  exit 1
fi

TARGET_DIR="$1"
SOURCE_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if [[ ! -d "$TARGET_DIR" ]]; then
  echo "Target directory does not exist: $TARGET_DIR" >&2
  exit 1
fi

mkdir -p "$TARGET_DIR/rules" "$TARGET_DIR/scripts" "$TARGET_DIR/.github/workflows"

cp "$SOURCE_ROOT/rules/unified-api.openapi.yaml" "$TARGET_DIR/rules/unified-api.openapi.yaml"
cp "$SOURCE_ROOT/rules/unified-api.human.md" "$TARGET_DIR/rules/unified-api.human.md"
cp "$SOURCE_ROOT/scripts/sync-api-docs.sh" "$TARGET_DIR/scripts/sync-api-docs.sh"
cp "$SOURCE_ROOT/.github/workflows/verify-api-doc-sync.yml" "$TARGET_DIR/.github/workflows/verify-api-doc-sync.yml"

chmod +x "$TARGET_DIR/scripts/sync-api-docs.sh"

LLM_RULES_FILE="$TARGET_DIR/LLM_RULES.md"
AGENTS_FILE="$TARGET_DIR/AGENTS.md"
CLAUDE_FILE="$TARGET_DIR/CLAUDE.md"
GEMINI_FILE="$TARGET_DIR/GEMINI.md"

cat > "$LLM_RULES_FILE" <<'EOF'
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
EOF

ensure_tool_entry() {
  local file="$1"
  local header="$2"
  local marker="## Unified API Rules Reference"

  if [[ ! -f "$file" ]]; then
    cat > "$file" <<EOF
# $header

$marker

Before starting any API/backend/frontend coding task, read:
- \`LLM_RULES.md\`
- \`rules/unified-api.openapi.yaml\`
- \`rules/unified-api.human.md\`

Enforcement:
- If conflict exists, OpenAPI wins.
- After API contract changes, run: \`bash scripts/sync-api-docs.sh\`
EOF
    return
  fi

  if ! grep -Fq "$marker" "$file"; then
    cat >> "$file" <<EOF

$marker

Before starting any API/backend/frontend coding task, read:
- \`LLM_RULES.md\`
- \`rules/unified-api.openapi.yaml\`
- \`rules/unified-api.human.md\`

Enforcement:
- If conflict exists, OpenAPI wins.
- After API contract changes, run: \`bash scripts/sync-api-docs.sh\`
EOF
  fi
}

ensure_tool_entry "$AGENTS_FILE" "AGENTS"
ensure_tool_entry "$CLAUDE_FILE" "CLAUDE"
ensure_tool_entry "$GEMINI_FILE" "GEMINI"

cat <<'EOF'
Bootstrapped API rule assets:
- rules/unified-api.openapi.yaml
- rules/unified-api.human.md
- scripts/sync-api-docs.sh
- .github/workflows/verify-api-doc-sync.yml
- LLM_RULES.md
- AGENTS.md / CLAUDE.md / GEMINI.md (created or patched with references)

Next steps:
1) Run:
   bash scripts/sync-api-docs.sh
2) Commit generated changes.
EOF
