#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

API_HOST="${API_HOST:-127.0.0.1:3005}"
DB_FILE="${DB_FILE:-/tmp/wasa-audit.db}"
SERVER_LOG="${SERVER_LOG:-/tmp/wasa-audit-server.log}"

BASE="http://${API_HOST}"

echo "== go test =="
cd "$ROOT_DIR"
go test ./...

echo "== golangci-lint =="
if [[ -x "$HOME/.local/bin/golangci-lint" ]]; then
  "$HOME/.local/bin/golangci-lint" run ./...
elif command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run ./...
else
  echo "golangci-lint not found. Install it, then re-run this script." >&2
  exit 2
fi

echo "== openapi lint (optional) =="
if [[ -x "$HOME/.npm-global/bin/lint-openapi" ]]; then
  "$HOME/.npm-global/bin/lint-openapi" "$ROOT_DIR/doc/api.yaml"
else
  echo "lint-openapi not found at $HOME/.npm-global/bin/lint-openapi (skipping)"
fi

echo "== webui lint/build =="
if [[ -f "$ROOT_DIR/webui/package.json" ]]; then
  if [[ ! -f "$ROOT_DIR/webui/node_modules/eslint/bin/eslint.js" ]]; then
    echo "webui/node_modules missing (or incomplete); running npm install..."
    (cd "$ROOT_DIR/webui" && npm install)
  fi
  (cd "$ROOT_DIR/webui" && npm run lint && npm run build-prod)
else
  echo "webui/package.json not found (skipping)"
fi

echo "== runtime API smoke test =="
python3 - <<'PY'
import base64, pathlib
b = base64.b64decode('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+XK6YAAAAASUVORK5CYII=')
pathlib.Path('/tmp/wasa_audit.png').write_bytes(b)
PY

rm -f "$DB_FILE"

CFG_WEB_API_HOST="$API_HOST" CFG_DB_FILENAME="$DB_FILE" CFG_DEBUG=true \
  go run ./cmd/webapi >"$SERVER_LOG" 2>&1 &
pid=$!
cleanup() { kill "$pid" 2>/dev/null || true; wait "$pid" 2>/dev/null || true; }
trap cleanup EXIT

for _ in $(seq 1 80); do
  if curl -sf "$BASE/liveness" >/dev/null; then break; fi
  sleep 0.1
done
curl -sf "$BASE/liveness" >/dev/null

# CORS preflight (required max-age=1)
cors_headers="$(curl -si -X OPTIONS "$BASE/session" \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Authorization,Content-Type" | tr -d '\r')"
echo "$cors_headers" | grep -qi '^Access-Control-Max-Age: 1$'

# Simplified login: {name} -> {identifier}
alice="$(curl -sf -X POST "$BASE/session" -H 'Content-Type: application/json' -d '{"name":"Alice"}' | python3 -c 'import sys,json; print(json.load(sys.stdin)["identifier"])')"
bob="$(curl -sf -X POST "$BASE/session" -H 'Content-Type: application/json' -d '{"name":"Bob"}' | python3 -c 'import sys,json; print(json.load(sys.stdin)["identifier"])')"

# User photo upload (Alice) + fetch (Bob)
curl -sf -o /dev/null -X PUT "$BASE/users/$alice/photo" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: image/png' \
  --data-binary @/tmp/wasa_audit.png
curl -sf -o /dev/null -X GET "$BASE/users/$alice/photo" \
  -H "Authorization: Bearer $bob" >/dev/null

# Direct message
curl -sf -o /dev/null -X POST "$BASE/users/$alice/chats/$bob/messages" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: application/json' \
  -d '{"body":"Hello Bob"}'

# Bob reads messages (triggers receipts) and reacts
mid="$(curl -sf -X GET "$BASE/users/$bob/chats/$alice/messages" -H "Authorization: Bearer $bob" \
  | python3 -c 'import sys,json; obj=json.load(sys.stdin); print(obj["messages"][0]["id"])')"
curl -sf -o /dev/null -X POST "$BASE/users/$bob/chats/$alice/messages/$mid/comments" \
  -H "Authorization: Bearer $bob" -H 'Content-Type: application/json' \
  -d '{"reaction":"😀"}'

# Alice sees status==2 and reaction present
curl -sf -X GET "$BASE/users/$alice/chats/$bob/messages" -H "Authorization: Bearer $alice" \
  | python3 -c 'import sys,json; obj=json.load(sys.stdin); m=obj["messages"][0]; assert m["status"]==2, m; assert any(r.get("reaction")=="😀" for r in m.get("reactions",[])), m'

# Forward + delete original
curl -sf -o /dev/null -X POST "$BASE/users/$alice/chats/$bob/messages/$mid/forward" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: application/json' \
  -d "{\"to\":\"$bob\"}"
curl -sf -o /dev/null -X DELETE "$BASE/users/$alice/chats/$bob/messages/$mid" \
  -H "Authorization: Bearer $alice"

# Create group (peer is g-<id>)
gid="$(curl -sf -X POST "$BASE/users/$alice/groups" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: application/json' \
  -d "{\"name\":\"TestGroup\",\"members\":[\"$bob\"]}" \
  | python3 -c 'import sys,json; print(json.load(sys.stdin)["group_id"])')"
peer="g-$gid"

# Group photo upload + fetch
curl -sf -o /dev/null -X PUT "$BASE/groups/$gid/photo" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: image/png' \
  --data-binary @/tmp/wasa_audit.png
curl -sf -o /dev/null -X GET "$BASE/groups/$gid/photo" \
  -H "Authorization: Bearer $bob" >/dev/null

# Group message
curl -sf -o /dev/null -X POST "$BASE/users/$alice/chats/$peer/messages" \
  -H "Authorization: Bearer $alice" -H 'Content-Type: application/json' \
  -d '{"body":"Hello group"}'

# Bob reads group (triggers receipts) and Alice must see status==2
curl -sf -X GET "$BASE/users/$bob/chats/$peer/messages" -H "Authorization: Bearer $bob" >/dev/null
curl -sf -X GET "$BASE/users/$alice/chats/$peer/messages" -H "Authorization: Bearer $alice" \
  | python3 -c 'import sys,json; obj=json.load(sys.stdin); m=obj["messages"][0]; assert m["status"]==2, m'

# Conversations list must be 200 JSON and include direct+group
curl -sf -X GET "$BASE/users/$alice/chats" -H "Authorization: Bearer $alice" \
  | python3 -c 'import json,sys; obj=json.load(sys.stdin); cs=obj.get("conversations",[]); assert any(c.get("isGroup") for c in cs), cs; assert any((not c.get("isGroup")) for c in cs), cs'

echo "AUDIT_OK"


