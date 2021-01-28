#!/bin/sh
set -e -u

DIR=$(cd $(dirname $0) && pwd)
cd "$DIR"

go generate
sleep 1
go test -v github.com/shamaton/msgpackgen/. -count=1