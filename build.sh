#!/usr/bin/env bash
# Builds one deployable zip per endpoint under dist/, for infra/lambda.tf to
# pick up. Each zip contains just a "bootstrap" binary (the provided.al2023
# custom runtime's expected entrypoint name). RSA signing keys are no longer
# baked into the zip; they're fetched at cold start from AWS Secrets Manager
# (see internal/bootstrap/app.go).
set -euo pipefail
cd "$(dirname "$0")"

DIST=dist
rm -rf "$DIST"
mkdir -p "$DIST"

ZIP_DIR="./zip_dir.py"

for dir in */; do
  name="${dir%/}"
  [ -f "$dir/main.go" ] || continue

  echo "building $name"
  work="$(mktemp -d)"
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o "$work/bootstrap" "./$name"

  python3 "$ZIP_DIR" "$work" "$DIST/$name.zip"
  rm -rf "$work"
done

echo "built $(ls "$DIST" | wc -l) function zips in $DIST/"
