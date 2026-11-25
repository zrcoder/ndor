package turtle

import (
	"io/fs"
	"math/rand/v2"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"

	"github.com/zrcoder/ndor/cmd/gui/internal/assets"
	"github.com/zrcoder/ndor/cmd/gui/internal/comp"
	turtleexamples "github.com/zrcoder/ndor/internal/examples/turtle"
	turtlelib "github.com/zrcoder/ndor/internal/turtle"
)

const (
	stepDelay      = 30 * time.Millisecond
	gopherFrameDur = 65 * time.Millisecond
	// stepSoundEvery plays a note every this many steps so dense moves do not
	// turn into a buzz.
	stepSoundEvery = 10
)

// stepNotes spans three octaves of the C major pentatonic, so random picks
// still sound consonant.
var stepNotes = [15]float64{
	261.63, 293.66, 329.63, 392.00, 440.00,
	523.25, 587.33, 659.25, 783.99, 880.00,
	1046.50, 1174.66, 1318.51, 1567.98, 1760.00,
}

type App struct {
	ws        *comp.Workspace
	prog      *turtlelib.Program
	gopher    []*ui.Bitmap // walk frames; frame 0 is the idle pose
	animStart time.Time
	stepTimer time.Time
	steps     int
}

// New builds a turtle app.
func New(fsys fs.FS) *App {
	a := &App{
		gopher: assets.SpriteSheet(fsys, "turtle_pen.png", 4, 64),
		prog:   turtlelib.New(),
	}
	a.ws = comp.NewWorkspace(fsys, comp.WorkspaceConfig{
		Code: `color "lightgreen"
width 3
for range 4 {
    step 10
    right 90
}
`,
		ExampleNames: turtleexamples.DefaultNames,
		Examples:     turtleexamples.Default,
		DocModule:    "turtle",
		Canvas:       a.canvas,
		OnRun:        a.runXGoCode,
	})
	return a
}

func (a *App) View(c *ui.Context) {
	a.ws.View(c)
}

// runXGoCode compiles the editor's xgo code and starts playback.
func (a *App) runXGoCode(c *ui.Context) {
	prog, err := turtlelib.Compile(a.ws.Editor.Buffer.String())
	if err != nil {
		a.prog = turtlelib.New()
		c.Toast(err.Error())
		c.Invalidate()
		return
	}

	a.prog = prog
	if a.prog.Done() {
		sound.Play(sound.Success)
	}
	a.steps = 0
	a.animStart = time.Now()
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

	if a.prog.Done() {
		sound.Play(sound.Success)
		return
	}
	a.steps++
	if a.steps%stepSoundEvery == 0 {
		sound.PlayNote(stepNotes[rand.IntN(len(stepNotes))])
	}
	p.AnimationFrame()
}

func (a *App) canvas(c *ui.Context) {
	ui.Box(c).WidthPercent(50).Fill().Background(ui.RGB(255, 255, 255)).DrawOver(func(p *ui.Painter, r ui.Rect) {
		a.tickStep(p)
		// Idle shows frame 0; playback cycles walk frames by time.
		frame := 0
		if !a.prog.Done() {
			frame = int(p.Now().Sub(a.animStart)/gopherFrameDur) % len(a.gopher)
		}
		a.prog.Draw(p, r, a.gopher[frame])
	})
}
