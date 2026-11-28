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
	result        *ndorlib.Drawing
	resultSVGData []byte
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
		Code: `color "lightgreen"
circle 0, 0, 300
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
		if painting || a.result == nil {
			return
		}
		ui.Row(c).Fill().Height(40).Children(func() {
			ui.Spacer(c)
			comp.IconButton(c, comp.IconButtonOption{Icon: icon.Must("download"), OnClick: a.saveImage})
		})
		ui.Box(c).Grow(1).FillWidth().DrawOver(func(p *ui.Painter, r ui.Rect) {
			r.X, r.W = r.X+10, max(r.W-20, 0)
			a.result.Render(p, r)
		})
	})
}

func (a *App) saveImage() {
	if a.resultSVGData == nil {
		return
	}

	m := locales.In(mygo.App.Locale())
	go func() {
		path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
			Title:       m.SaveImage(),
			DefaultPath: "ndor.svg",
			Filters:     []mygo.FileFilter{{Name: m.SvgImage(), Extensions: []string{"svg"}}},
		})
		if err != nil {
			a.notify(err)
			return
		}
		if path == "" {
			return
		}
		err = os.WriteFile(path, a.resultSVGData, 0o644)
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
	code := a.ws.Editor.Buffer.String()

	d, lineErr := ndorlib.Compile(code)
	if lineErr != nil {
		a.result = nil
		a.resultSVGData = nil
		a.ws.MakeError(lineErr.Number, lineErr.Msg)
		c.Invalidate()
		return
	}
	a.ws.ClearError()

	a.result = d
	a.resultSVGData = d.SVG()
	a.paintUntil = time.Now().Add(comp.PaintingDuration)
	sound.Play(sound.Swipe)
	c.Invalidate()
}
