package comp

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/zrcoder/ndor/internal"
	"github.com/zrcoder/ndor/pkg"
)

const (
	turtleStepDur  = 160 * time.Millisecond // 每步基准时长
	turtleMaxTotal = 8 * time.Second        // 整段动画上限
	turtleCell     = 40                     // 网格格数（turtle 像素）
)

// Turtle replays a turtle plan with a small animation: the turtle draws
// step by step across a coordinate grid, then rests on the final frame.
func Turtle(c *ui.Context, plan *pkg.TurtlePlan, runId uint64) {
	t := c.Theme()
	n := len(plan.Steps)
	total := min(time.Duration(n)*turtleStepDur, turtleMaxTotal)

	box := ui.Box(c.Key("turtle")).Fill()
	started := ui.Local(box, runId, func() time.Time { return time.Now() })

	box.Draw(func(p *ui.Painter, r ui.Rect) {
		scale := min(r.W/float32(plan.W), r.H/float32(plan.H))
		ox := r.X + (r.W-float32(plan.W)*scale)/2
		oy := r.Y + (r.H-float32(plan.H)*scale)/2
		p.Transform(ui.Translation(ox+scale*float32(plan.W)/2, oy+scale*float32(plan.H)/2), 0, 0, func() {
			p.Transform(ui.Scaling(scale, -scale), 0, 0, func() {
				drawGrid(p, t, float32(plan.W), float32(plan.H))

				now := p.Now()
				elapsed := now.Sub(*started)
				done := elapsed >= total

				prog := float64(elapsed) / float64(total)
				if prog > 1 {
					prog = 1
				}
				iEnd := min(int(prog*float64(n)), n)
				frac := prog*float64(n) - float64(iEnd)
				if frac > 1 {
					frac = 1
				}

				drawSegments(p, plan.Steps, iEnd, float32(frac))

				var x, y, h float64
				if n > 0 && iEnd > 0 {
					last := plan.Steps[iEnd-1]
					if iEnd < n {
						seg := plan.Steps[iEnd]
						x = last.X2 + float64(frac)*(seg.X2-last.X2)
						y = last.Y2 + float64(frac)*(seg.Y2-last.Y2)
						h = seg.H
					} else {
						x, y, h = last.X2, last.Y2, last.H
					}
				}
				drawTurtle(p, float32(x), float32(y), float32(h), t.Accent)

				if !done {
					p.AnimationFrame()
				}
			})
		})
	})
}

func drawGrid(p *ui.Painter, t *ui.Theme, w, h float32) {
	grid := t.Border
	grid.A = 120
	grid5 := t.Border
	grid5.A = 200
	axis := t.TextMuted

	for x := -math.Ceil(float64(w/2/turtleCell)) * turtleCell; x <= math.Ceil(float64(w/2/turtleCell))*turtleCell; x += turtleCell {
		c := grid
		if int(x)%(turtleCell*5) == 0 {
			c = grid5
		}
		if x == 0 {
			c = axis
		}
		p.Line(float32(x), -h/2, float32(x), h/2, 1, c)
	}
	for y := -math.Ceil(float64(h/2/turtleCell)) * turtleCell; y <= math.Ceil(float64(h/2/turtleCell))*turtleCell; y += turtleCell {
		c := grid
		if int(y)%(turtleCell*5) == 0 {
			c = grid5
		}
		if y == 0 {
			c = axis
		}
		p.Line(-w/2, float32(y), w/2, float32(y), 1, c)
	}
}

func drawSegments(p *ui.Painter, steps []internal.TurtleStep, done int, frac float32) {
	for i, s := range steps {
		if i > done {
			break
		}
		f := 1.0
		if i == done {
			f = float64(frac)
		}
		if f <= 0 {
			continue
		}
		x2 := s.X1 + float64(f)*(s.X2-s.X1)
		y2 := s.Y1 + float64(f)*(s.Y2-s.Y1)
		c := ui.RGBA(s.R, s.G, s.B, float32(s.A)/255)
		p.Line(float32(s.X1), float32(s.Y1), float32(x2), float32(y2), float32(s.LineW), c)
	}
}

func drawTurtle(p *ui.Painter, x, y, h float32, c ui.Color) {
	p.Transform(ui.Translation(x, y), 0, 0, func() {
		p.Transform(ui.Rotation(h), 0, 0, func() {
			tr := (&ui.Path{}).MoveTo(0, -9).LineTo(-6, 7).LineTo(6, 7).Close()
			p.FillPath(tr, c)
		})
	})
}
