let _codeEditor = null;
let _errorMarker = null;

window.addEventListener("load", function () {
  const container = document.getElementById("codeArea");
  if (!container) {
    console.error("Editor container not found");
    return;
  }

  container.style.height = "100%";
  container.style.width = "100%";

  console.log("codeArea dimensions:", container.offsetWidth, container.offsetHeight);

  if (typeof CodeMirror === "undefined") {
    console.error(
      "CodeMirror is not loaded. Make sure CDN links are included.",
    );
    return;
  }

  console.log("CodeMirror loaded");
  initEditor();
});

function initEditor() {
  try {
    _codeEditor = CodeMirror(document.getElementById("codeArea"), {
      mode: "text/x-csrc",
      theme: "monokai",
      lineNumbers: true,
      lineWrapping: true,
      scrollbarStyle: "null",
      autofocus: true,
      lineHeight: 24,
      fontSize: 16,
      smartIndent: false,
    });

    _codeEditor.setCursor({ line: 3000, ch: 0 });

    window.addEventListener("resize", function () {
      if (_codeEditor) {
        _codeEditor.refresh();
      }
    });

    console.log("CodeMirror Editor initialized");
  } catch (e) {
    console.error("Failed to initialize CodeMirror Editor:", e);
  }
}

window.MarkErrorLine = function (number) {
  if (!_codeEditor) return;

  if (_errorMarker) {
    _codeEditor.removeLineClass(_errorMarker, "background", "editorLineErr");
  }

  _errorMarker = _codeEditor.addLineClass(
    number - 1,
    "background",
    "editorLineErr",
  );
};

window.SetCode = function (s) {
  if (!_codeEditor) return;
  _codeEditor.setValue(s);
  _codeEditor.setCursor({ line: 3000, ch: 0 });
  _codeEditor.focus();
};

window.GetCode = function () {
  if (!_codeEditor) return "";
  return _codeEditor.getValue();
};
