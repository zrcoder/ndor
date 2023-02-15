require.config({ paths: { vs: './js/lib/monaco-editor/vs' } })

let codeEditor
require(['vs/editor/editor.main'], function () {
    codeEditor = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'c', // go+
        fontSize: 20,
    })
    codeEditor.focus()
    codeEditor.onDidChangeModelContent(() => {
        clearErrorLineMarks()
    })
})

let decorations = []

function clearErrorLineMarks() {
    codeEditor.removeDecorations(decorations)
}

function markErrorLine(number) {
    decorations = codeEditor.deltaDecorations([],
        [
            {
                range: new monaco.Range(number, 0, number, 0),
                options: {
                    isWholeLine: true,
                    inlineClassName: "editorLineErr"
                },
            },
        ]
    )
}

function getCode() {
    return codeEditor.getValue()
}

function setCode(s) {
    codeEditor.setValue(s)
    codeEditor.setPosition({ column: 1000, lineNumber: 3000 })
    codeEditor.focus()
}
