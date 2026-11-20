package turtle

import (
	"io/fs"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"

	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
	"github.com/zrcoder/ndor/cmd/gui/internal/comp"
	turtleexamples "github.com/zrcoder/ndor/internal/examples/turtle"
	turtlelib "github.com/zrcoder/ndor/internal/turtle"
)

const stepDelay = 80 * time.Millisecond

type App struct {
	prog      *turtlelib.Program
	editor    *editor.State
	help      *comp.Help
	helpImg   *ui.Bitmap
	runImg    *ui.Bitmap
	gopher    *ui.Bitmap
	stepTimer time.Time
}

// New builds a turtle app.
func New(fsys fs.FS) *App {
	a := &App{
		helpImg: assets.Bitmap(fsys, "teacher_labeled.png"),
		runImg:  assets.Bitmap(fsys, "gopher_labeled.png"),
		gopher:  assets.Sprite(fsys, "gopher_pen.svg", 50),
		prog:    turtlelib.New(),
	}

	code := `color "lightgreen"
width 3
for range 4 {
    step 10
    right 90
}
`
	a.editor = &editor.State{Buffer: editor.NewBuffer(code)}
	a.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))

	a.help = comp.NewHelp(comp.HelpConfig{
		DialogImg:    assets.Bitmap(fsys, "teacher.png"),
		CloseIcon:    assets.SVG(fsys, "close.svg"),
		ExampleNames: turtleexamples.DefaultNames,
		Examples:     turtleexamples.Default,
		OnPickExample: func(code string) {
			a.editor.Buffer.Set(code)
		},
		LoadDoc: assets.DocLoader(fsys, "turtle"),
	})
	return a
}

// runXGoCode compiles the editor's xgo code and starts playback.
func (a *App) runXGoCode(c *ui.Context) {
	prog, err := turtlelib.Compile(a.editor.Buffer.String())
	if err != nil {
		a.prog = turtlelib.New()
		c.Toast(err.Error())
		c.Invalidate()
		return
	}

	a.prog = prog
	a.stepTimer = time.Now().Add(-stepDelay)
	if !a.prog.Done() {
		c.Invalidate()
	}
}

// tickStep applies one queued operation per frame.
func (a *App) tickStep(p *ui.Painter) {
	if a.prog.Done() {
		return
	}
	if !a.stepTimer.IsZero() && p.Now().Sub(a.stepTimer) < stepDelay {
		p.AnimationFrame()
		return
	}

	a.prog.Advance()
	a.stepTimer = p.Now()

	if !a.prog.Done() {
		p.AnimationFrame()
	}
}

func (a *App) View(c *ui.Context) {
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		a.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		a.runXGoCode(c)
	}

	a.help.View(c)

	ui.Row(c.Key("main")).Fill().Children(func() {
		ui.Box(c).WidthPercent(50).Fill().Background(ui.RGB(255, 255, 255)).DrawOver(func(p *ui.Painter, r ui.Rect) {
			a.tickStep(p)
			a.prog.Draw(p, r, a.gopher)
		})

		ui.Column(c).Fill().Children(func() {
			editor.View(c, a.editor).Fill()
		})
	})

	comp.FloatButtons(c, a.helpImg, a.runImg, a.help.Open, func() { a.runXGoCode(c) })
}
