#!/usr/bin/env ash
#
# Builds and archives project with resources
# in a golang:alpine Docker container.

set -euo pipefail

ROOT_DIR="/satelit-index"
TARGET_DIR="docker/satelit-index"
ARTIFACTS_DIR="$ROOT_DIR/$TARGET_DIR"

make_install() {
  echo "Building project"

  CGO_ENABLED=0 \
    go build -a -o satelit-index
  mv satelit-index "$ARTIFACTS_DIR"
}

install_tools() {
  echo "Installing tools"

  GO111MODULE=on \
    GOBIN="$ARTIFACTS_DIR" \
    CGO_ENABLED=0 \
      go get -u github.com/pressly/goose/cmd/goose
}

copy_resources() {
  echo "Copying resources"

  mkdir -p "$ARTIFACTS_DIR/migrations" "$ARTIFACTS_DIR/config"
  cp config/*.yml "$ARTIFACTS_DIR/config"
  cp sql/schema/*.sql "$ARTIFACTS_DIR/migrations"
  cp docker/scripts/entry.sh "$ARTIFACTS_DIR"
}

archive() {
  echo "Packing artifacts"

  apk add tar
  find "$TARGET_DIR/" -type f -o -type l -o -type d \
    | sed s,^"$TARGET_DIR/",, \
    | tar -czf satelit-index.tar.gz \
    --no-recursion -C "$TARGET_DIR/" -T -
}

main() {
  if [[ -d "$ARTIFACTS_DIR" ]]; then
    rm -rf "$ARTIFACTS_DIR"
  fi

  mkdir "$ARTIFACTS_DIR"

  make_install
  install_tools
  copy_resources
  archive
}

main "$@"
