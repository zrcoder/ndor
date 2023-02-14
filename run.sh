#! /usr/bin/env bash

wwwDir=./static
monacoEditorDir=${wwwDir}/monaco-editor
serverDir=./server

cp -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ${wwwDir}
GOOS=js GOARCH=wasm go build -o ./static/niudour.wasm .
go run ${serverDir}