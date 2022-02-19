require.config({ paths: { vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.31.1/min/vs' } })
var codeEditor
require(['vs/editor/editor.main'], function () {
    const codeArea = document.getElementById('codeArea')
    codeEditor = monaco.editor.create(codeArea, {
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
    codeEditor.setPosition({column: 1000, lineNumber: 3000})
    codeEditor.focus()
}