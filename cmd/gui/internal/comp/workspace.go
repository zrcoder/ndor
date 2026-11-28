package comp

import (
	"io/fs"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"

	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
)

// WorkspaceConfig configures the shared workspace shell.
type WorkspaceConfig struct {
	Code         string
	ExampleNames []string
	Examples     map[string]string
	DocModule    string
	Canvas       func(c *ui.Context)
	OnRun        func(c *ui.Context)
}

// Workspace lays out the canvas, code editor, help dialog and floating buttons.
type Workspace struct {
	Editor   *editor.State
	cfg      WorkspaceConfig
	help     *Help
	settings *Settings
	copyBtn  *CopyButton
	helpImg  *ui.Bitmap
	runImg   *ui.Bitmap
	errLine  int
	errMsg   string
	errVer   uint64
}

// NewWorkspace creates a workspace.
func NewWorkspace(fsys fs.FS, cfg WorkspaceConfig) *Workspace {
	w := &Workspace{
		cfg:     cfg,
		helpImg: assets.Bitmap(fsys, "teacher_labeled.png"),
		runImg:  assets.Bitmap(fsys, "gopher_labeled.png"),
	}
	w.Editor = &editor.State{Buffer: editor.NewBuffer(cfg.Code)}
	w.Editor.Buffer.SetLanguage(editor.LanguageNamed("c"))

	w.help = NewHelp(HelpConfig{
		DialogImg:     assets.Bitmap(fsys, "teacher.png"),
		ExampleNames:  cfg.ExampleNames,
		Examples:      cfg.Examples,
		OnPickExample: w.Editor.Buffer.Set,
		LoadDoc:       assets.DocLoader(fsys, cfg.DocModule),
	})

	w.copyBtn = NewCopyButton(func(*ui.Context) string { return w.Editor.Buffer.String() })
	w.settings = NewSettings()
	return w
}

// MakeError marks line, from 1, and shows msg above it.
func (w *Workspace) MakeError(line int, msg string) {
	w.errLine, w.errMsg, w.errVer = line, msg, w.Editor.Buffer.Version()
}

// ClearError removes the marks.
func (w *Workspace) ClearError() {
	w.errLine, w.errMsg = 0, ""
	w.Editor.LineBackgrounds, w.Editor.GutterMarks = nil, nil
	w.Editor.SetDiagnostics(w.Editor.Buffer.Version(), nil)
}

// errorLine returns the error's line from 0, or -1.
func (w *Workspace) errorLine() int {
	line := w.errLine - 1
	if w.errMsg == "" || line < 0 || line >= w.Editor.Buffer.LineCount() {
		return -1
	}
	return line
}

// markError highlights the error line, dropping it once the code changes.
func (w *Workspace) markError(t *ui.Theme) {
	buf, v := w.Editor.Buffer, w.Editor.Buffer.Version()
	if w.errLine != 0 && v != w.errVer {
		w.ClearError()
		return
	}
	if line := w.errorLine(); line >= 0 {
		start := buf.Offset(editor.Pos{Line: line})
		w.Editor.LineBackgrounds = []editor.LineBackground{{First: line, Last: line, Color: t.Danger.Alpha(0.12)}}
		w.Editor.GutterMarks = []editor.GutterMark{{Line: line, Color: t.Danger}}
		w.Editor.SetDiagnostics(v, []editor.Diagnostic{{
			Start: start, End: start + len(buf.Line(line)), Severity: editor.SeverityError, Message: w.errMsg,
		}})
	}
}

// viewError floats the message above its line.
func (w *Workspace) viewError(c *ui.Context) {
	line := w.errorLine()
	if line < 0 {
		return
	}
	t := c.Theme()
	top := max(float32(w.Editor.LineTop(line)-w.Editor.ScrollY())+t.Space(5), 0)
	ui.Row(c.Key("error")).Absolute().Left(t.Space(14)).Top(top).Gap(t.Space(1)).
		Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Space(0.75)).Background(t.Danger).
		Shadow(0, 2, 6, 0, ui.RGBA(0, 0, 0, 0.2)).IgnoreInput(true).Children(func() {
		white := ui.RGB(255, 255, 255)
		ui.Icon(c, icon.Must("circle-alert")).Size(12, 12).TextColor(white)
		ui.Text(c, w.errMsg).FontSize(11).TextColor(white).SingleLine()
	})
}

// View renders the workspace.
func (w *Workspace) View(c *ui.Context) {
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		w.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		w.cfg.OnRun(c)
	}

	w.help.View(c)
	w.settings.View(c)
	w.markError(c.Theme())

	ui.Row(c.Key("main")).Fill().Children(func() {
		w.cfg.Canvas(c)
		ui.Column(c).Fill().Children(func() {
			ui.Row(c).Fill().Height(36).Children(func() {
				ui.Spacer(c)
				w.copyBtn.View(c)
				IconButton(c, IconButtonOption{Icon: icon.Must("settings"), OnClick: w.settings.Open})
			})
			ui.Box(c.Key("code")).Fill().Children(func() {
				editor.View(c, w.Editor).Fill()
				w.viewError(c)
			})
		})
	})

	FloatButtons(c, w.helpImg, w.runImg, w.help.Open, func() { w.cfg.OnRun(c) })
}
