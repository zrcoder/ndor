package ndor

import (
	"io/fs"
	"os"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	ndorlib "github.com/zrcoder/ndor"
	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
	"github.com/zrcoder/ndor/cmd/gui/internal/comp"
	examples "github.com/zrcoder/ndor/internal/examples/ndor"
)

type App struct {
	pendingToast  string
	resultImg     *ui.Bitmap
	resultPNGData []byte
	paintUntil    time.Time
	editor        *editor.State
	helpBtnImg    *ui.Bitmap
	runBtnImg     *ui.Bitmap
	paintImg      *ui.Bitmap
	codeImg       *ui.Bitmap
	help          *comp.Help
	copyButton    *comp.CopyButton
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

	a.copyButton = comp.NewCopyButton(comp.CopyButtonConfig{
		CopyIcon:  assets.SVG(fsys, "copy.svg"),
		CheckIcon: assets.SVG(fsys, "check.svg"),
		GetCode:   func() string { return a.editor.Buffer.String() },
	})

	return a
}

func (a *App) View(c *ui.Context) {
	if a.pendingToast != "" {
		c.Toast(a.pendingToast)
		a.pendingToast = ""
	}
	// Handle shortcuts before anything else so a run shows this frame.
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		a.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		a.runCode(c)
	}

	a.help.View(c)

	painting := comp.PaintingActive(c, a.paintUntil)
	if painting {
		c.After(time.Until(a.paintUntil))
	}

	ui.Row(c.Key("main")).Fill().Children(func() {
		a.pictureArea(c, painting)
		a.codeArea(c)
	})

	comp.FloatButtons(c, a.helpBtnImg, a.runBtnImg, a.help.Open, func() { a.runCode(c) })

	if painting {
		comp.Painting(c, a.paintImg, a.codeImg)
	}
}

func (a *App) pictureArea(c *ui.Context, painting bool) {
	ui.Column(c).WidthPercent(50).Background(ui.RGB(255, 255, 255)).Fill().Children(func() {
		if painting || a.resultImg == nil {
			return
		}
		ui.Image(c, a.resultImg).MaxWidthPercent(100).MaxHeightPercent(100).ContextMenu(func(m *ui.Menu) {
			if m.Item("Save").Chosen() {
				a.saveImage()
			}
		})
	})
}

func (a *App) saveImage() {
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
		if err != nil {
			a.notify("Save failed: " + err.Error())
			return
		}
		if path == "" {
			return
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			a.notify("Save failed: " + err.Error())
			return
		}
		a.notify("Saved")
	}()
}

// notify shows a toast from any goroutine. Toasts only register during a
// frame, so the toast waits for the next View.
func (a *App) notify(msg string) {
	mygo.RunOnMain(func() {
		a.pendingToast = msg
		if win := mygo.FocusedWindow(); win != nil {
			win.Invalidate()
		}
	})
}

func (a *App) codeArea(c *ui.Context) {
	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Fill().Height(36).Children(func() {
			ui.Spacer(c)
			a.copyButton.View(c)
		})
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
		a.resultImg = nil
		a.resultPNGData = nil
		c.Toast(lineErr.Error())
		c.Invalidate()
		return
	}

	a.resultPNGData = pngData
	bitmap, err := ui.DecodeBitmap(pngData)
	if err != nil {
		a.resultImg = nil
		a.resultPNGData = nil
		c.Toast("Failed to decode image: " + err.Error())
		c.Invalidate()
		return
	}

	a.resultImg = bitmap
	a.paintUntil = time.Now().Add(comp.PaintingDuration)
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
