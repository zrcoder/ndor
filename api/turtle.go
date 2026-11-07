package api

import (
	"github.com/zrcoder/ndor/internal"
)

// Turtle places the turtle at (x, y) in canvas pixels measured from the
// canvas center with y up. Default is the canvas center, facing up.
func Turtle(x, y float64) {
	internal.Turtle.Jump(x, y)
}

// Forward moves the turtle forward by d canvas pixels, drawing a line while
// the pen is down.
func Forward(d float64) {
	internal.Turtle.Forward(d)
}

// Back moves the turtle backward by d canvas pixels.
func Back(d float64) {
	internal.Turtle.Back(d)
}

// Left turns the turtle left by deg degrees.
func Left(deg float64) {
	internal.Turtle.Left(deg)
}

// Right turns the turtle right by deg degrees.
func Right(deg float64) {
	internal.Turtle.Right(deg)
}

// Angle sets the turtle's absolute heading: 0 points up, clockwise positive.
func Angle(deg float64) {
	internal.Turtle.Angle(deg)
}

// Jump moves the turtle to (x, y) without drawing.
func Jump(x, y float64) {
	internal.Turtle.Jump(x, y)
}

// PenUp lifts the pen so the turtle moves without drawing.
func PenUp() {
	internal.Turtle.Pen(false)
}

// PenDown puts the pen down so the turtle draws while moving.
func PenDown() {
	internal.Turtle.Pen(true)
}

// Home moves the turtle back to the canvas center, facing up, pen down.
func Home() {
	internal.Turtle.Home()
}
