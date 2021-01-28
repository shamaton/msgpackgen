#!/usr/bin/env bash

DIR=$(dirname "$0")
CURRENT=$(cd "$DIR" && pwd)
PJ=$(cd "$CURRENT" && cd ../ && pwd)
cd "${PJ}" || exit 1

find . -name "define*_test.go" -maxdepth 1 | sed 's/_test.go//' | xargs -I{} mv {}_test.go {}.go;
go generate
find . -name "define*.go" -maxdepth 1 | sed 's/.go//' | xargs -I{} mv {}.go {}_test.go;
go test -v -count=1 \
  --coverpkg=github.com/shamaton/msgpackgen/... \
  --coverprofile=coverage.coverprofile \
  --covermode=atomic ./...
git checkout resolver_test.go