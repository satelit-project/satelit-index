#!/usr/bin/env sh
#
# Entry point for running built project
# in an alpine Docker container.

set -euo pipefail

wait_db() {
  local retries=5
  while [[ "$retries" -gt "0" ]]; do
    set +e
    ./goose postgres "$PG_DB_URL" version
    local status="$?"
    set -e

    if [[ "$status" -eq "0" ]]; then
      echo "Database available"
      return
    fi

    retries=$(( retries - 1 ))
    echo "Database is not available. Sleeping..."
    sleep 10s
  done

  exit 1
}

main() {
  echo "Waiting for database"
  wait_db

  echo "Running migrations"
  ./goose -dir migrations \
    postgres "$PG_DB_URL" \
    up

  echo "Running service"
  ST_LOG=prod \
    DO_SPACES_KEY="$DO_SPACES_KEY" \
    DO_SPACES_SECRET="$DO_SPACES_SECRET" \
    DO_SPACES_HOST="$DO_SPACES_HOST" \
    DO_BUCKET="$DO_BUCKET" \
    PG_DB_URL="$PG_DB_URL" \
    exec ./satelit-index
}

main "$@"
