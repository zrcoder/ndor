require.config({ paths: { vs: './js/lib/monaco-editor/vs' } })

let _codeEditor = null
let _decorations = []

require(['vs/editor/editor.main'], function () {
    _codeEditor = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'c', // go+
        fontSize: 20,
        wordWrapColumn: 45,
    })
    _codeEditor.focus()
    _codeEditor.onDidChangeModelContent(() => {
        _codeEditor.removeDecorations(_decorations)
        _decorations = []
    })
})

function MarkErrorLine(number) {
    _decorations = _codeEditor.deltaDecorations([], [{
        range: new monaco.Range(number, 0, number, 0),
        options: {
            isWholeLine: true,
            inlineClassName: "editorLineErr"
        }
    }])
}

function SetCode(s) {
    _codeEditor.setValue(s)
    _codeEditor.setPosition({ column: 0, lineNumber: 3000 })
    _codeEditor.focus()
}

function GetCode() {
    return _codeEditor.getValue()
}