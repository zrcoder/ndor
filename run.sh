#! /usr/bin/env bash

wwwDir=./static
monacoEditorDir=${wwwDir}/monaco-editor
serverDir=./server

function prepareMonacoEditor() {
  if [ -d ${monacoEditorDir} ]; then
    return
  fi
  echo "begin to download monaco-editor"
  mkdir -p ${monacoEditorDir}
  curl https://registry.npmjs.org/monaco-editor/-/monaco-editor-0.29.1.tgz -o ${monacoEditorDir}/manoco-editor.tgz
  tar -zxf ${monacoEditorDir}/manoco-editor.tgz -C ${monacoEditorDir}
  rm ${monacoEditorDir}/manoco-editor.tgz
}

function build() {
    cp -f "$(go env GOROOT)/misc/wasm/wasm_exec.js" ${wwwDir}
    GOOS=js GOARCH=wasm go build -o ./static/niudour.wasm .
}

function run() {
    go run ${serverDir}
}

prepareMonacoEditor
build
run

