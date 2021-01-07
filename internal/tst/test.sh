#!/bin/sh
set -e -u

DIR=$(cd $(dirname $0) && pwd)
cd "$DIR"

go generate
go test -v github.com/shamaton/msgpackgen/... -count=1