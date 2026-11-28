// Turtle drawing commands used by xgo code.
// Commands collect in CmdQueue and are replayed frame by frame by the app.

package turtle

import "github.com/zrcoder/ndor/internal/parse"

// cmd is one micro-operation produced by the user's xgo code.
type cmd struct {
	Kind  string // "step" | "turn" | "penUp" | "penDown" | "setColor" | "setWidth" | "clear" | "home"
	Val   float64
	Color [4]uint8
}

// CmdQueue holds all commands produced by the user's code, in order.
var CmdQueue []cmd

func push(kind string, val float64) {
	CmdQueue = append(CmdQueue, cmd{Kind: kind, Val: val})
}

// Step moves the turtle d steps (user units, 1 = 1 grid cell);
// a negative d moves it backward.
func Step(d float64) { push("step", d) }

// Left turns the turtle left by d degrees.
func Left(d float64) { push("turn", d) }

// Right turns the turtle right by d degrees.
func Right(d float64) { push("turn", -d) }

// PenColor sets the pen color. Accepts a color name string (e.g. "red", "#ff0000")
// or 3-4 numeric args (r, g, b[, a]); alpha defaults to 255.
func PenColor(args ...any) {
	c, ok := parse.Color(args...)
	if !ok {
		return
	}
	CmdQueue = append(CmdQueue, cmd{Kind: "setColor", Color: [4]uint8{c.R, c.G, c.B, c.A}})
}

// PenWidth sets the pen width.
func PenWidth(w float64) { push("setWidth", w) }

// Clear clears the drawing and resets the turtle.
func Clear() {
	CmdQueue = nil
	CmdQueue = append(CmdQueue, cmd{Kind: "clear"})
}

// Home resets the turtle position and angle without clearing the drawing.
func Home() { push("home", 0) }

// Pen lifts or lowers the pen. Accepts "up" or "down" (default).
func Pen(args ...string) {
	if len(args) > 0 && args[0] == "up" {
		push("penUp", 0)
		return
	}
	push("penDown", 0)
}
