#!/usr/bin/env bash
CWD=$(pwd)
cd $GOPATH/src/github.com/gen0cide/gscript/cmd/gscript
go build -o $GOPATH/bin/gscript
cd $CWD
