require.config({
  paths: {
    vs: "./js/lib/monaco-editor/vs",
  },
  waitSeconds: 60,
});

let _codeEditor = null;
let _decorations = [];

window.addEventListener("load", function () {
  const container = document.getElementById("codeArea");
  if (!container) {
    console.error("Editor container not found");
    return;
  }

  container.style.height = "100%";
  container.style.width = "100%";

  require(["vs/editor/editor.main"], function () {
    initEditor();
  });
});

function initEditor() {
  try {
    _codeEditor = monaco.editor.create(document.getElementById("codeArea"), {
      language: "c",
      theme: "vs-dark",
      fontSize: 16,
      wordWrap: "on",
      value: "\n",
      minimap: {
        enabled: false,
      },
      scrollbar: {
        vertical: "hidden",
        horizontal: "hidden",
      },
      automaticLayout: true,
      overviewRulerLanes: 0,
      hideCursorInOverviewRuler: true,
    });

    _codeEditor.focus();
    _codeEditor.setPosition({ column: 0, lineNumber: 3000 });

    _codeEditor.onDidChangeModelContent(() => {
      if (_decorations) {
        _codeEditor.removeDecorations(_decorations);
        _decorations = [];
      }
    });

    window.addEventListener("resize", function () {
      if (_codeEditor) {
        _codeEditor.layout();
      }
    });
  } catch (e) {
    console.error("Failed to initialize Monaco Editor:", e);
  }
}

window.MarkErrorLine = function (number) {
  if (!_codeEditor) return;
  _decorations = _codeEditor.deltaDecorations(
    [],
    [
      {
        range: new monaco.Range(number, 0, number, 0),
        options: {
          isWholeLine: true,
          inlineClassName: "editorLineErr",
        },
      },
    ]
  );
};

window.SetCode = function (s) {
  if (!_codeEditor) return;
  _codeEditor.setValue(s);
  _codeEditor.setPosition({ column: 0, lineNumber: 3000 });
  _codeEditor.focus();
};

window.GetCode = function () {
  if (!_codeEditor) return "";
  return _codeEditor.getValue();
};
