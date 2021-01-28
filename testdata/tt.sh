#!/bin/sh
set -e -u

DIR=$(cd $(dirname $0) && pwd)
cd "$DIR"

go generate
sleep 1
cd ..
go test -v -count=1 --coverpkg=github.com/shamaton/msgpackgen/... --coverprofile=coverage.coverprofile --covermode=atomic ./...

git checkout testdata/resolver.go