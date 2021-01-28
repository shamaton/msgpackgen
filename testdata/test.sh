#!/usr/bin/env bash

DIR=$(dirname "$0")
CURRENT=$(cd "$DIR" && pwd)
PJ=$(cd "$CURRENT" && cd ../ && pwd)
cd "${PJ}" || exit 1

pwd
exit 0

find . -name "define*_test.go" -exec rename -v 's/_test//i' {} \;
go generate
find . -name "define*.go" -exec rename -v 's/\.go$/_test\.go/i' {} \;
go test -v .