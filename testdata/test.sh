#!/usr/bin/env bash
set -e

DIR=$(dirname "$0")
CURRENT=$(cd "$DIR" && pwd)
PJ=$(cd "$CURRENT" && cd ../ && pwd)
cd "${PJ}" || exit 1

TMP_GENERATED=$(mktemp)
cp msgpack.msgpackgen_test.go "${TMP_GENERATED}"
cleanup() {
  cp "${TMP_GENERATED}" msgpack.msgpackgen_test.go
  rm -f "${TMP_GENERATED}"
  rm -f msgpack.msgpackgen.go
}
trap cleanup EXIT

go generate
go test -v -count=1 \
  --coverpkg=github.com/shamaton/msgpackgen/... \
  --coverprofile=coverage.coverprofile \
  --covermode=atomic ./...
