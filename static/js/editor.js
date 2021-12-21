require.config({ paths: { vs: './monaco-editor/package/min/vs' } })
var editModel
require(['vs/editor/editor.main'], function () {
    editModel = monaco.editor.create(document.getElementById('codeArea'), {
        language: 'coffeescript',
        fontSize: 24,
    })
})

function getCode() {
    return editModel.getValue()
}