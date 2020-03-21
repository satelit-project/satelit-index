#!/usr/bin/env sh
#
# Entry point for running built project
# in an alpine Docker container.

set -euo pipefail

assert_env() {
  if declare -p "$1" &>/dev/null; then
    return
  fi

  echo -e "\033[0;31mVariable '""$1""'is not set.\033[0m"
  exit 1
}

main() {
  assert_env "DO_SPACES_KEY"
  assert_env "DO_SPACES_SECRET"
  assert_env "DO_SPACES_HOST"
  assert_env "DO_BUCKET"
  assert_env "PG_DB_URL"

  # run migrations
  ./goose -dir sql \
    postgres "$PG_DB_URL" \
    up

  # run service
  ST_LOG=prod \
    DO_SPACES_KEY="$DO_SPACES_KEY" \
    DO_SPACES_SECRET="$DO_SPACES_SECRET" \
    DO_SPACES_HOST="$DO_SPACES_HOST" \
    DO_BUCKET="$DO_BUCKET" \
    PG_DB_URL="$PG_DB_URL" \
    exec ./satelit-index
}

main "$@"
