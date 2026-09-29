#!/bin/sh
set -eu
: "${DATABASE_URL:?DATABASE_URL is required}"
PSQL_BIN=${PSQL_BIN:-psql}
for file in migrations/*.sql; do
  if command -v "$PSQL_BIN" >/dev/null 2>&1; then
    "$PSQL_BIN" "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$file"
  else
    docker compose exec -T db psql "$DATABASE_URL" -v ON_ERROR_STOP=1 < "$file"
  fi
done
