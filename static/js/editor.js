require.config({ paths: { vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.35.0/min/vs' } })
var codeEditor
require(['vs/editor/editor.main'], function () {
    codeEditor = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'coffeescript',
        fontSize: 24,
    })
    codeEditor.focus()
})

function getCode() {
    return codeEditor.getValue()
}

function setCode(s) {
    codeEditor.setValue(s)
    codeEditor.setPosition({ column: 1000, lineNumber: 3000 })
    codeEditor.focus()
}