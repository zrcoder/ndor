#! /usr/bin/env bash

wwwDir=./static
serverDir=./server

cp -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ${wwwDir}
GOOS=js GOARCH=wasm go build -o ${wwwDir}/niudour.wasm .
go run ${serverDir}