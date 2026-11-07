package internal

import (
	icolor "image/color"
	"math"
)

// TurtleStep is one pen-down move of the turtle: a segment in canvas
// pixels, drawn in pen color and width. H is the heading in degrees while
// drawing, 0 pointing up and clockwise positive, for the replay marker.
type TurtleStep struct {
	X1, Y1, X2, Y2 float64
	R, G, B, A     uint8
	LineW          float64
	H              float64
}

// Pen is the turtle's pen: color and width, as Color and Linew set them.
var (
	PenColor = icolor.RGBA{A: 0xff}
	PenW     = 1.0
)

// turtleState is the turtle's state: x and y in canvas pixels measured from
// the canvas center up as positive, heading in degrees (0 up, clockwise
// positive) and pen, with the canvas size for pixel conversions.
type turtleState struct {
	x, y    float64
	heading float64
	pen     bool
	w, h    float64
}

var (
	Turtle      turtleState
	TurtleSteps []TurtleStep
)

func ResetTurtle(width, height int) {
	Turtle = turtleState{pen: true, w: float64(width), h: float64(height)}
	TurtleSteps = nil
}

func (t *turtleState) Forward(d float64) {
	nx := t.x + d*math.Sin(math.Pi*t.heading/180)
	ny := t.y + d*math.Cos(math.Pi*t.heading/180)
	if t.pen {
		TurtleSteps = append(TurtleSteps, TurtleStep{
			X1: t.ix(t.x), Y1: t.iy(t.y),
			X2: t.ix(nx), Y2: t.iy(ny),
			R: PenColor.R, G: PenColor.G, B: PenColor.B, A: PenColor.A,
			LineW: PenW, H: t.heading,
		})
	}
	t.x, t.y = nx, ny
}

func (t *turtleState) Back(d float64)    { t.Forward(-d) }
func (t *turtleState) Left(deg float64)  { t.heading -= deg }
func (t *turtleState) Right(deg float64) { t.heading += deg }
func (t *turtleState) Angle(deg float64) { t.heading = deg }
func (t *turtleState) Jump(x, y float64) { t.x, t.y = x, y }
func (t *turtleState) Home()             { t.x, t.y = 0, 0; t.heading = 0; t.pen = true }
func (t *turtleState) Pen(down bool)     { t.pen = down }

// End reports where the turtle rests, in canvas pixels and heading degrees.
func (t *turtleState) End() (x, y, h float64) {
	return t.ix(t.x), t.iy(t.y), t.heading
}

// ix and iy map a turtle coordinate, measured from the canvas center with y
// up, to a canvas pixel.
func (t *turtleState) ix(x float64) float64 { return t.w/2 + x }
func (t *turtleState) iy(y float64) float64 { return t.h/2 - y }
