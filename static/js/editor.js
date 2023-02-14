require.config({ paths: { vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.35.0/min/vs' } })
var codeEditor
require(['vs/editor/editor.main'], function () {
    codeEditor = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'c', // go+
        fontSize: 24,
        wordWrap: 'on',
    })
    codeEditor.focus()
    codeEditor.onDidChangeModelContent(() => {
        clearErrorLineMarks()
    })
})

let decorations = []

function markErrorLine(number) {
    decorations = codeEditor.deltaDecorations(
        [],
        [
            {
                range: new monaco.Range(number, 0, number, 0),
                options: {
                    isWholeLine: true,
                    inlineClassName: "editorLineErr"
                },
            },
        ])
}

function clearErrorLineMarks() {
    codeEditor.removeDecorations(decorations)
}

function getCode() {
    return codeEditor.getValue()
}

function setCode(s) {
    codeEditor.setValue(s)
    codeEditor.setPosition({ column: 1000, lineNumber: 3000 })
    codeEditor.focus()
}
