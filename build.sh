#!/usr/bin/env bash
# Builds one deployable zip per endpoint under dist/, for infra/lambda.tf to
# pick up. Each zip contains just a "bootstrap" binary (the provided.al2023
# custom runtime's expected entrypoint name) plus this service's RSA keys,
# since lexi-users-ms signs tokens and needs both.
set -euo pipefail
cd "$(dirname "$0")"

DIST=dist
rm -rf "$DIST"
mkdir -p "$DIST"

KEYS_SRC="../../microservicios/lexi-users-ms/src/keys"
ZIP_DIR="./zip_dir.py"

for dir in */; do
  name="${dir%/}"
  [ -f "$dir/main.go" ] || continue

  echo "building $name"
  work="$(mktemp -d)"
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o "$work/bootstrap" "./$name"
  cp "$KEYS_SRC/private.pem" "$KEYS_SRC/public.pem" "$work/"

  python3 "$ZIP_DIR" "$work" "$DIST/$name.zip"
  rm -rf "$work"
done

echo "built $(ls "$DIST" | wc -l) function zips in $DIST/"
