#!/bin/sh
set -e -u

DIR=$(cd $(dirname $0) && pwd)
cd "$DIR"

go generate
sleep 1
go test -v github.com/shamaton/msgpackgen/... -count=1
sleep 1
git checkout resolver.msgpackgen.go