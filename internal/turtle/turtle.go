package turtle

import "math"

// pt is a point in math space: origin at center, +Y up.
type pt struct{ X, Y float32 }

// seg is a continuous stroke with its own color and width.
type seg struct {
	pts   []pt
	color [3]uint8
	width float32
}

type state struct {
	x, y, angle float32 // degrees; 0 points right
	down        bool
	penColor    [3]uint8
	penWidth    float32
	segments    []seg
	cur         *seg
}

const defaultWidth = 2

func (t *state) turn(d float32) { t.angle += d }

func (t *state) penDown() {
	t.down = true
	if t.cur == nil {
		t.cur = &seg{color: t.penColor, width: t.penWidth}
		t.cur.pts = append(t.cur.pts, pt{t.x, t.y})
	}
}

func (t *state) penUp() {
	t.down = false
	if t.cur != nil && len(t.cur.pts) > 0 {
		t.segments = append(t.segments, *t.cur)
		t.cur = nil
	}
}

func (t *state) home() {
	t.x, t.y, t.angle = 0, 0, 0
	if t.down {
		t.penUp()
		t.penDown()
	}
}

func (t *state) clear() {
	t.x, t.y, t.angle = 0, 0, 0
	t.down = true
	t.penColor = [3]uint8{20, 20, 20}
	t.penWidth = defaultWidth
	t.segments = nil
	t.cur = nil
	t.penDown()
}

func (t *state) setColor(r, g, b uint8) {
	t.penColor = [3]uint8{r, g, b}
	t.splitStroke()
}

func (t *state) setWidth(w float32) {
	t.penWidth = w
	t.splitStroke()
}

// splitStroke ends the current stroke so later changes keep old points intact.
func (t *state) splitStroke() {
	if t.cur != nil && len(t.cur.pts) > 1 {
		t.segments = append(t.segments, *t.cur)
		t.cur = nil
	}
	if t.down && t.cur == nil {
		t.penDown()
	}
}

// step moves d (negative = backward), recording the point when the pen is down.
func (t *state) step(d float32) {
	rad := t.angle * math.Pi / 180
	nx := t.x + float32(math.Cos(float64(rad)))*d
	ny := t.y + float32(math.Sin(float64(rad)))*d
	if t.down {
		if t.cur == nil {
			t.penDown()
		}
		t.cur.pts = append(t.cur.pts, pt{nx, ny})
	}
	t.x, t.y = nx, ny
}
