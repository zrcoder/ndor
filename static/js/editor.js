let decorations = null;
let editor = null;

require.config({
    paths: {
        'vs': 'https://cdn.jsdelivr.net/npm/monaco-editor@0.27.0/min/vs'
    }
});

require(['vs/editor/editor.main'], function() {
    decorations = []
    editor = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'c', // go+
        theme: 'vs-dark',
        fontSize: 16,
        wordWrap: 'on',
        minimap: {
            enabled: false
        },
        scrollbar: {
            vertical: "hidden",
            horizontal: "hidden"
        },
        automaticLayout: true,
        overviewRulerLanes: 0,
        hideCursorInOverviewRuler: true,
    })

    editor.onDidChangeModelContent(() => {
        editor.removeDecorations(decorations);
        decorations = [];
    })
});


function getCode() {
    return editor.getValue();
}

function markErrorLine(number) {
    decorations = editor.deltaDecorations([], [{
        range: new monaco.Range(number, 0, number, 0),
        options: {
            isWholeLine: true,
            inlineClassName: "editorLineErr"
        }
    }])
}

function setCode(s) {
    editor.setValue(s);
    editor.setPosition({
        column: 0,
        lineNumber: 3000
    })
    editor.focus();
}