package ndor

import (
	"io/fs"
	"os"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	ndorlib "github.com/zrcoder/ndor"
	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
	"github.com/zrcoder/ndor/cmd/gui/internal/comp"
	"github.com/zrcoder/ndor/internal/examples/ndor"
)

type App struct {
	resultImg     *ui.Bitmap
	resultPNGData []byte
	errorMsg      string
	errorLine     int
	editor        *editor.State
	helpBtnImg    *ui.Bitmap
	runBtnImg     *ui.Bitmap
	paintImg      *ui.Bitmap
	codeImg       *ui.Bitmap
	help          *comp.Help
	toolbar       *comp.Toolbar
}

func New(fsys fs.FS) *App {
	a := &App{
		helpBtnImg: assets.Bitmap(fsys, "teacher_labeled.png"),
		runBtnImg:  assets.Bitmap(fsys, "gopher_labeled.png"),
		paintImg:   assets.Bitmap(fsys, "paint.png"),
		codeImg:    assets.Bitmap(fsys, "code.png"),
	}
	code := `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
`
	a.editor = &editor.State{Buffer: editor.NewBuffer(code)}
	a.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))

	a.help = comp.NewHelp(comp.HelpConfig{
		DialogImg:    assets.Bitmap(fsys, "teacher.png"),
		CloseIcon:    assets.SVG(fsys, "close.svg"),
		ExampleNames: examples.DefaultNames,
		Examples:     examples.Default,
		OnPickExample: func(code string) {
			a.editor.Buffer.Set(code)
		},
		LoadDoc: assets.DocLoader(fsys, "ndor"),
	})

	a.toolbar = comp.NewToolbar(comp.ToolbarConfig{
		CopyIcon:  assets.SVG(fsys, "copy.svg"),
		CheckIcon: assets.SVG(fsys, "check.svg"),
		SaveIcon:  assets.SVG(fsys, "save.svg"),
		CloseIcon: assets.SVG(fsys, "close.svg"),
		GetCode:   func() string { return a.editor.Buffer.String() },
		OnCopy:    func(string) {},
		OnSave:    func() { a.saveImage(nil) },
	})
	return a
}

func (a *App) View(c *ui.Context) {
	// Handle shortcuts before the toast viewport so a run shows this frame.
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		a.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		a.runCode(c)
	}

	if a.errorMsg != "" {
		comp.ToastErrorAtLine(c, a.errorLine, a.errorMsg, func() {
			a.errorMsg = ""
			a.errorLine = 0
		})
	}

	a.help.View(c)

	painting := comp.PaintingViewport(c, a.paintImg, a.codeImg)

	ui.Row(c.Key("main")).Fill().Children(func() {
		a.pictureArea(c, painting)
		a.codeArea(c)
	})

	ui.Column(c.Key("help-float")).Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		comp.ImageButtonSimple(c, a.helpBtnImg, func() {
			a.help.Open()
		})
	})

	ui.Column(c.Key("run-float")).Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		comp.ImageButtonSimple(c, a.runBtnImg, func() {
			a.runCode(c)
		})
	})
}

func (a *App) pictureArea(c *ui.Context, painting bool) {
	ui.Column(c).WidthPercent(50).Fill().Children(func() {
		if painting {
			return
		}
		if a.resultImg == nil {
			return
		}
		ui.Row(c).Fill().Height(40).Children(func() {
			ui.Spacer(c)
			if a.toolbar != nil {
				a.toolbar.SaveButton(c)
			}
		})
		ui.Image(c, a.resultImg).MaxWidthPercent(100).MaxHeightPercent(100)
	})
}

func (a *App) saveImage(_ *ui.Context) {
	if a.resultPNGData == nil {
		return
	}
	data := a.resultPNGData
	go func() {
		path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
			Title:       "Save Image",
			DefaultPath: "ndor.png",
			Filters:     []mygo.FileFilter{{Name: "PNG Image", Extensions: []string{"png"}}},
		})
		if err != nil || path == "" {
			return
		}
		os.WriteFile(path, data, 0o644)
	}()
}

func (a *App) codeArea(c *ui.Context) {
	ui.Column(c).Fill().Children(func() {
		if a.toolbar != nil {
			a.toolbar.View(c)
		}
		editor.View(c, a.editor).Fill()
	})
}

func (a *App) runCode(c *ui.Context) {
	width, height := getPictureSize(c)
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}
	code := a.editor.Buffer.String()

	pngData, lineErr := ndorlib.Gen(width, height, code)
	if lineErr != nil {
		a.errorMsg = lineErr.Msg
		a.errorLine = lineErr.Number
		a.resultImg = nil
		a.resultPNGData = nil
		c.Invalidate()
		return
	}

	a.errorMsg = ""
	a.errorLine = 0
	a.resultPNGData = pngData
	bitmap, err := ui.DecodeBitmap(pngData)
	if err != nil {
		a.errorMsg = "Failed to decode image: " + err.Error()
		a.errorLine = 0
		a.resultImg = nil
		a.resultPNGData = nil
		c.Invalidate()
		return
	}

	a.resultImg = bitmap

	comp.ShowPainting(c)
	c.Invalidate()
}

func getPictureSize(c *ui.Context) (int, int) {
	win := mygo.FocusedWindow()
	if win != nil {
		bounds := win.ContentBounds()
		return bounds.Width / 2, bounds.Height - 40 - 72
	}
	return 800, 600
}
