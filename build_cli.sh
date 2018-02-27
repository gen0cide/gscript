#!/usr/local/bin/bash
set -x
CWD=$(pwd)
cd $GOPATH/src/github.com/gen0cide/gscript/cmd/gscript
go-bindata -pkg compiler -nomemcopy -o $GOPATH/src/github.com/gen0cide/gscript/compiler/bindata.go -prefix '../..' ../../templates/...
go build -o $GOPATH/bin/gscript
cd $CWD
