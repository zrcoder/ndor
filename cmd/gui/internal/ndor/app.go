package ndor

import (
	"io/fs"
	"os"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	ndorlib "github.com/zrcoder/ndor"
	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
	"github.com/zrcoder/ndor/cmd/gui/internal/comp"
	"github.com/zrcoder/ndor/cmd/gui/internal/locales"
	examples "github.com/zrcoder/ndor/internal/examples/ndor"
)

type App struct {
	ws            *comp.Workspace
	resultImg     *ui.Bitmap
	resultPNGData []byte
	paintUntil    time.Time
	painting      bool
	paintImg      *ui.Bitmap
	codeImg       *ui.Bitmap
	notice        comp.Notice
}

func New(fsys fs.FS) *App {
	a := &App{
		paintImg: assets.Bitmap(fsys, "paint.png"),
		codeImg:  assets.Bitmap(fsys, "code.png"),
	}
	a.ws = comp.NewWorkspace(fsys, comp.WorkspaceConfig{
		Code: `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
`,
		ExampleNames: examples.DefaultNames,
		Examples:     examples.Default,
		DocModule:    "ndor",
		Canvas:       a.canvas,
		OnRun:        a.runCode,
	})
	return a
}

func (a *App) View(c *ui.Context) {
	a.ws.View(c)

	painting := comp.PaintingActive(c, a.paintUntil)
	if a.painting && !painting {
		sound.Play(sound.Success)
	}
	a.painting = painting
	if painting {
		comp.Painting(c, a.paintImg, a.codeImg)
	}
	a.notice.View(c)
}

func (a *App) canvas(c *ui.Context) {
	painting := comp.PaintingActive(c, a.paintUntil)
	if painting {
		c.After(time.Until(a.paintUntil))
	}
	ui.Themed(c, ui.LightTheme()).WidthPercent(50).Fill().Children(func() {
		if painting || a.resultImg == nil {
			return
		}
		ui.Row(c).Fill().Height(40).Children(func() {
			ui.Spacer(c)
			comp.IconButton(c, comp.IconButtonOption{Icon: icon.Must("download"), OnClick: a.saveImage})
		})
		ui.Image(c, a.resultImg).MaxWidthPercent(100).MaxHeightPercent(100)
	})
}

func (a *App) saveImage() {
	if a.resultPNGData == nil {
		return
	}

	m := locales.In(mygo.App.Locale())
	go func() {
		path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
			Title:       m.SaveImage(),
			DefaultPath: "ndor.png",
			Filters:     []mygo.FileFilter{{Name: m.PngImage(), Extensions: []string{"png"}}},
		})
		if err != nil {
			a.notify(err)
			return
		}
		if path == "" {
			return
		}
		err = os.WriteFile(path, a.resultPNGData, 0o644)
		a.notify(err)
	}()
}

// notify shows the save result, from the goroutine that saved the file.
func (a *App) notify(err error) {
	mygo.RunOnMain(func() {
		m := locales.In(mygo.App.Locale())
		if err != nil {
			a.notice.Show(feedback.ToneDanger, m.SaveFailed(), err.Error(), 3*time.Second)
		} else {
			a.notice.Show(feedback.ToneSuccess, m.Saved(), "", 2*time.Second)
		}
	})
}

func (a *App) runCode(c *ui.Context) {
	width, height := getPictureSize()
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}
	code := a.ws.Editor.Buffer.String()

	pngData, lineErr := ndorlib.Gen(width, height, code)
	if lineErr != nil {
		a.resultImg = nil
		a.resultPNGData = nil
		a.ws.MakeError(lineErr.Number, lineErr.Msg)
		c.Invalidate()
		return
	}
	a.ws.ClearError()

	a.resultPNGData = pngData
	bitmap, err := ui.DecodeBitmap(pngData)
	if err != nil {
		a.resultImg = nil
		a.resultPNGData = nil
		a.notice.Show(feedback.ToneDanger, locales.In(c.Locale()).DecodeFailed(err.Error()), "", 3*time.Second)
		return
	}

	a.resultImg = bitmap
	a.paintUntil = time.Now().Add(comp.PaintingDuration)
	sound.Play(sound.Swipe)
	c.Invalidate()
}

func getPictureSize() (int, int) {
	win := mygo.FocusedWindow()
	if win != nil {
		bounds := win.ContentBounds()
		return bounds.Width / 2, bounds.Height - 40 - 72
	}
	return 800, 600
}
