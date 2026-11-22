package comp

import (
	"io/fs"

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
	Editor  *editor.State
	cfg     WorkspaceConfig
	help    *Help
	copyBtn *CopyButton
	helpImg *ui.Bitmap
	runImg  *ui.Bitmap
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
		CloseIcon:     assets.SVG(fsys, "close.svg"),
		ExampleNames:  cfg.ExampleNames,
		Examples:      cfg.Examples,
		OnPickExample: w.Editor.Buffer.Set,
		LoadDoc:       assets.DocLoader(fsys, cfg.DocModule),
	})

	w.copyBtn = NewCopyButton(CopyButtonConfig{
		CopyIcon:  assets.SVG(fsys, "copy.svg"),
		CheckIcon: assets.SVG(fsys, "check.svg"),
		GetCode:   func() string { return w.Editor.Buffer.String() },
	})
	return w
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

	ui.Row(c.Key("main")).Fill().Children(func() {
		w.cfg.Canvas(c)
		ui.Column(c).Fill().Children(func() {
			ui.Row(c).Fill().Height(36).Children(func() {
				ui.Spacer(c)
				w.copyBtn.View(c)
			})
			editor.View(c, w.Editor).Fill()
		})
	})

	FloatButtons(c, w.helpImg, w.runImg, w.help.Open, func() { w.cfg.OnRun(c) })
}
