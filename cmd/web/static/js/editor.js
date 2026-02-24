let _codeEditor = null;

function tryInitCodeMirror() {
  const container = document.getElementById("codeArea");
  if (!container) {
    return false;
  }

  if (typeof CodeMirror === "undefined") {
    return false;
  }

  if (_codeEditor) {
    const wrapper = _codeEditor.getWrapperElement();
    if (wrapper && container.contains(wrapper)) {
      return true;
    }
    console.log("DOM replaced, reinitializing CodeMirror...");
    _codeEditor = null;
  }

  container.style.height = "100%";
  container.style.width = "100%";

  console.log("Initializing CodeMirror");

  try {
    _codeEditor = CodeMirror(container, {
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
    return true;
  } catch (e) {
    console.error("Failed to initialize CodeMirror Editor:", e);
    return false;
  }
}

function startWatching() {
  // Wait for CodeMirror to be available
  function waitForCodeMirror() {
    if (typeof CodeMirror === "undefined") {
      setTimeout(waitForCodeMirror, 100);
      return;
    }
    initWatcher();
  }

  function initWatcher() {
    const observer = new MutationObserver(function () {
      tryInitCodeMirror();
    });

    // Watch the entire body for any changes
    observer.observe(document.body, { childList: true, subtree: true });
    console.log("Started watching body for changes");

    // Also try immediately in case everything is already ready
    tryInitCodeMirror();
  }

  waitForCodeMirror();
}

window.addEventListener("load", function () {
  startWatching();
});

document.addEventListener("visibilitychange", function () {
  if (document.visibilityState === "visible") {
    tryInitCodeMirror();
  }
});

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
