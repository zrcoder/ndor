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
	x, y, angle float32 // degrees; 0 points right, default faces up (90)
	down        bool
	penColor    [3]uint8
	penWidth    float32
	segments    []seg
	cur         []pt // in-progress stroke points; style comes from penColor/penWidth
}

const defaultWidth = 2

func (t *state) turn(d float32) { t.angle += d }

// startCur begins a new in-progress stroke at the current position when the pen is down.
func (t *state) startCur() {
	if t.down && t.cur == nil {
		t.cur = []pt{{t.x, t.y}}
	}
}

// flushCur archives the current stroke (if it has drawn points) using the current pen style.
func (t *state) flushCur() {
	if len(t.cur) > 1 {
		t.segments = append(t.segments, seg{pts: t.cur, color: t.penColor, width: t.penWidth})
	}
	t.cur = nil
}

func (t *state) penDown() {
	t.down = true
	t.startCur()
}

func (t *state) penUp() {
	t.down = false
	t.flushCur()
}

func (t *state) home() {
	t.x, t.y, t.angle = 0, 0, 90
	t.flushCur()
	t.startCur()
}

func (t *state) clear() {
	t.x, t.y, t.angle = 0, 0, 90
	t.penColor = [3]uint8{20, 20, 20}
	t.penWidth = defaultWidth
	t.segments = nil
	t.cur = nil
	t.penDown()
}

// setColor archives the current stroke (old color) before switching to the new one.
func (t *state) setColor(r, g, b uint8) {
	t.flushCur()
	t.penColor = [3]uint8{r, g, b}
	t.startCur()
}

// setWidth archives the current stroke (old width) before switching to the new one.
func (t *state) setWidth(w float32) {
	t.flushCur()
	t.penWidth = w
	t.startCur()
}

// step moves d (negative = backward), recording the point when the pen is down.
func (t *state) step(d float32) {
	rad := t.angle * math.Pi / 180
	nx := t.x + float32(math.Cos(float64(rad)))*d
	ny := t.y + float32(math.Sin(float64(rad)))*d
	t.startCur()
	if t.down {
		t.cur = append(t.cur, pt{nx, ny})
	}
	t.x, t.y = nx, ny
}
