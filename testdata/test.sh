#!/usr/bin/env bash
set -e

DIR=$(dirname "$0")
CURRENT=$(cd "$DIR" && pwd)
PJ=$(cd "$CURRENT" && cd ../ && pwd)
cd "${PJ}" || exit 1

TMP_RESOLVER=$(mktemp)
cp resolver_test.go "${TMP_RESOLVER}"
cleanup() {
  cp "${TMP_RESOLVER}" resolver_test.go
  rm -f "${TMP_RESOLVER}"
  rm -f resolver.msgpackgen.go
}
trap cleanup EXIT

go generate
go test -v -count=1 \
  --coverpkg=github.com/shamaton/msgpackgen/... \
  --coverprofile=coverage.coverprofile \
  --covermode=atomic ./...
