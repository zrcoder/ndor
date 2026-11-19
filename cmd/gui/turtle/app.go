package turtle

import (
	"math"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"

	"github.com/zrcoder/ndor/cmd/gui/comp"
	"github.com/zrcoder/ndor/cmd/gui/turtle/api"
)

const (
	cell               = 10
	cw                 = 560
	ch                 = 440
	stepDelay          = 80 * time.Millisecond
	turnDelay          = 150 * time.Millisecond
	stepPixelsPerFrame = 4
)

// stepOp is one micro-operation produced by the user's xgo code.
type stepOp struct {
	kind  string // "step" | "turn" | "penUp" | "penDown" | "setColor" | "setWidth" | "clear" | "home"
	val   float32
	color [3]uint8
}

// Config holds the ndor-owned assets the turtle window reuses.
type Config struct {
	HelpButtonImg *ui.Bitmap
	RunButtonImg  *ui.Bitmap
	HelpDialogImg *ui.Bitmap
	CloseIcon     *ui.SVG
	LoadDoc       func(isEnglish bool, tab int) string
}

type App struct {
	win       *mygo.Window
	t         turtle
	editor    *editor.State
	help      *comp.Help
	helpImg   *ui.Bitmap
	runImg    *ui.Bitmap
	errorMsg  string
	stepQueue []stepOp
	stepIdx   int
	stepTimer time.Time
}

// New builds a turtle app; call SetWindow once its mygo window exists.
func New(cfg Config) (*App, error) {
	a := &App{helpImg: cfg.HelpButtonImg, runImg: cfg.RunButtonImg}
	a.t.clear()
	if err := initGopher(); err != nil {
		return nil, err
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

	names := make([]string, 0, len(examples))
	byName := make(map[string]string, len(examples))
	for _, ex := range examples {
		names = append(names, ex.name)
		byName[ex.name] = ex.code
	}
	a.help = comp.NewHelp(comp.HelpConfig{
		DialogImg:    cfg.HelpDialogImg,
		CloseIcon:    cfg.CloseIcon,
		ExampleNames: names,
		Examples:     byName,
		OnPickExample: func(code string) {
			a.editor.Buffer.Set(code)
		},
		LoadDoc: cfg.LoadDoc,
	})
	return a, nil
}

func (a *App) SetWindow(w *mygo.Window) { a.win = w }

func (a *App) pushStep(op stepOp) {
	a.stepQueue = append(a.stepQueue, op)
}

// runXGoCode runs the editor's xgo code and queues drawing operations.
func (a *App) runXGoCode(c *ui.Context) {
	a.stepQueue = nil
	a.stepIdx = 0
	a.stepTimer = time.Time{}
	a.t.clear()

	api.CmdQueue = nil

	if err := runCode(a.editor.Buffer.String()); err != nil {
		a.errorMsg = err.Error()
		c.Invalidate()
		return
	}
	a.errorMsg = ""

	for _, cmd := range api.CmdQueue {
		switch cmd.Kind {
		case "step":
			dist := float32(cmd.Val) * cell
			n := max(int(math.Ceil(float64(dist)/float64(stepPixelsPerFrame))), 1)
			step := dist / float32(n)
			for range n {
				a.pushStep(stepOp{kind: "step", val: step})
			}
		case "turn":
			a.pushStep(stepOp{kind: "turn", val: float32(cmd.Val)})
		case "penUp":
			a.pushStep(stepOp{kind: "penUp"})
		case "penDown":
			a.pushStep(stepOp{kind: "penDown"})
		case "setColor":
			a.pushStep(stepOp{kind: "setColor", color: cmd.Color})
		case "setWidth":
			a.pushStep(stepOp{kind: "setWidth", val: float32(cmd.Val)})
		case "clear":
			a.pushStep(stepOp{kind: "clear"})
		case "home":
			a.pushStep(stepOp{kind: "home"})
		}
	}

	if len(a.stepQueue) > 0 {
		a.stepTimer = time.Now().Add(-stepDelay)
		a.win.Invalidate()
	}
}

// tickStep runs one queued operation per frame.
func (a *App) tickStep(p *ui.Painter) {
	if a.stepIdx >= len(a.stepQueue) {
		return
	}
	if !a.stepTimer.IsZero() && p.Now().Sub(a.stepTimer) < stepDelay {
		p.AnimationFrame()
		return
	}

	op := a.stepQueue[a.stepIdx]
	a.stepIdx++
	a.stepTimer = p.Now()

	// Mutate state in the draw callback, avoiding win.Update.
	switch op.kind {
	case "step":
		a.t.step(op.val)
	case "turn":
		a.t.turn(op.val)
	case "penUp":
		a.t.penUp()
	case "penDown":
		a.t.penDown()
	case "setColor":
		a.t.setColor(op.color[0], op.color[1], op.color[2])
	case "setWidth":
		a.t.setWidth(op.val)
	case "clear":
		a.t.clear()
	case "home":
		a.t.home()
	}

	if a.stepIdx < len(a.stepQueue) {
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

	if a.errorMsg != "" {
		comp.Toast(c, a.errorMsg, func() { a.errorMsg = "" })
	}

	a.help.View(c)

	ui.Row(c.Key("main")).Fill().Children(func() {
		ui.Box(c).WidthPercent(50).Fill().Background(ui.RGB(255, 255, 255)).DrawOver(func(p *ui.Painter, r ui.Rect) {
			a.tickStep(p)
			a.draw(p, r)
		})

		ui.Column(c).Fill().Children(func() {
			editor.View(c, a.editor).Fill()
		})
	})

	ui.Column(c.Key("help-float")).Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		comp.ImageButtonSimple(c, a.helpImg, func() { a.help.Open() })
	})
	ui.Column(c.Key("run-float")).Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		comp.ImageButtonSimple(c, a.runImg, func() { a.runXGoCode(c) })
	})
}

func (a *App) draw(p *ui.Painter, r ui.Rect) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2

	// Grid lines darken with distance from the origin; the axes are darkest.
	gridLight := ui.RGB(240, 240, 240)
	gridMid := ui.RGB(220, 220, 220)
	gridCenter := ui.RGB(150, 150, 150)

	strokeGrid := func(x0, y0, x1, y1 float32, col ui.Color) {
		ln := ui.Path{}
		ln.MoveTo(x0, y0).LineTo(x1, y1)
		p.StrokePath(&ln, 1, col)
	}

	nx := int(cx / cell)
	for k := -nx; k <= nx; k++ {
		x := cx + float32(k)*cell
		if x < r.X || x > r.X+r.W {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(x, r.Y, x, r.Y+r.H, col)
	}

	ny := int(cy / cell)
	for k := -ny; k <= ny; k++ {
		y := cy + float32(k)*cell
		if y < r.Y || y > r.Y+r.H {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(r.X, y, r.X+r.W, y, col)
	}

	sx := func(q pt) float32 { return cx + q.X }
	sy := func(q pt) float32 { return cy - q.Y }

	strokeSeg := func(s seg) {
		if len(s.pts) < 2 {
			return
		}
		ln := ui.Path{}
		ln.MoveTo(sx(s.pts[0]), sy(s.pts[0]))
		for _, q := range s.pts[1:] {
			ln.LineTo(sx(q), sy(q))
		}
		p.StrokePath(&ln, s.width, ui.RGB(s.color[0], s.color[1], s.color[2]))
	}
	for _, s := range a.t.segments {
		strokeSeg(s)
	}
	if a.t.cur != nil {
		strokeSeg(*a.t.cur)
	}

	px, py := sx(pt{a.t.x, a.t.y}), sy(pt{a.t.x, a.t.y})
	drawGopher(p, px, py, a.t.angle)
}
