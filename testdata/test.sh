#!/usr/bin/env bash
set -e

DIR=$(dirname "$0")
CURRENT=$(cd "$DIR" && pwd)
PJ=$(cd "$CURRENT" && cd ../ && pwd)
cd "${PJ}" || exit 1

go generate
go test -v -count=1 \
  --coverpkg=github.com/shamaton/msgpackgen/... \
  --coverprofile=coverage.coverprofile \
  --covermode=atomic ./...
git checkout resolver_test.go
rm -f resolver.msgpackgen.go