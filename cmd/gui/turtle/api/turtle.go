// Package api provides the turtle drawing commands used by xgo code.
// Commands collect in CmdQueue and are replayed frame by frame by the app.
package api

import (
	"errors"

	icolor "image/color"

	"golang.org/x/image/colornames"
)

// cmd is one micro-operation produced by the user's xgo code.
type cmd struct {
	Kind  string // "step" | "turn" | "penUp" | "penDown" | "setColor" | "setWidth" | "clear" | "home"
	Val   float64
	Color [3]uint8
}

// CmdQueue holds all commands produced by the user's code, in order.
var CmdQueue []cmd

// Step moves the turtle d steps (user units, 1 = 1 grid cell);
// a negative d moves it backward.
func Step(d float64) {
	CmdQueue = append(CmdQueue, cmd{Kind: "step", Val: d})
}

// Left turns the turtle left by d degrees.
func Left(d float64) {
	CmdQueue = append(CmdQueue, cmd{Kind: "turn", Val: d})
}

// Right turns the turtle right by d degrees.
func Right(d float64) {
	CmdQueue = append(CmdQueue, cmd{Kind: "turn", Val: -d})
}

// PenColor sets the pen color. Accepts a color name string (e.g. "red", "#ff0000")
// or 3-4 numeric args (r, g, b[, a]).
func PenColor(args ...any) {
	c, ok := parseColor(args)
	if !ok {
		return
	}
	turtleColor = c
	CmdQueue = append(CmdQueue, cmd{Kind: "setColor", Color: c})
}

// PenWidth sets the pen width.
func PenWidth(w float64) {
	turtleWidth = w
	CmdQueue = append(CmdQueue, cmd{Kind: "setWidth", Val: w})
}

// Clear clears the drawing and resets the turtle.
func Clear() {
	turtleX, turtleY, turtleAngle = 0, 0, 0
	turtleDown = true
	turtleColor = [3]uint8{20, 20, 20}
	turtleWidth = 2.0
	CmdQueue = nil
	CmdQueue = append(CmdQueue, cmd{Kind: "clear"})
}

// Home resets the turtle position and angle without clearing the drawing.
func Home() {
	turtleX, turtleY, turtleAngle = 0, 0, 0
	CmdQueue = append(CmdQueue, cmd{Kind: "home"})
}

// Pen lifts or lowers the pen. Accepts "up" or "down" (default).
func Pen(args ...string) {
	if len(args) > 0 && args[0] == "up" {
		turtleDown = false
		CmdQueue = append(CmdQueue, cmd{Kind: "penUp"})
		return
	}
	turtleDown = true
	CmdQueue = append(CmdQueue, cmd{Kind: "penDown"})
}

func parseColor(args []any) ([3]uint8, bool) {
	var c icolor.RGBA
	switch len(args) {
	case 1:
		s, ok := args[0].(string)
		if !ok {
			return [3]uint8{}, false
		}
		if col, found := colornames.Map[s]; found {
			c = col
		} else if col, found := parseHexColor(s); found {
			c = col
		} else {
			return [3]uint8{}, false
		}
	case 3, 4:
		floats, err := toFloats(args)
		if err != nil {
			return [3]uint8{}, false
		}
		var a uint8 = 0xff
		if len(args) == 4 {
			a = uint8(floats[3])
		}
		c = icolor.RGBA{R: uint8(floats[0]), G: uint8(floats[1]), B: uint8(floats[2]), A: a}
	default:
		return [3]uint8{}, false
	}
	return [3]uint8{c.R, c.G, c.B}, true
}

func parseHexColor(s string) (icolor.RGBA, bool) {
	if s == "" || s[0] != '#' {
		return icolor.RGBA{}, false
	}
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		default:
			return 0
		}
	}
	c := icolor.RGBA{A: 0xff}
	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		return icolor.RGBA{}, false
	}
	return c, true
}

func toFloats(pars []any) ([]float64, error) {
	res := make([]float64, len(pars))
	for i, s := range pars {
		n, err := parseFloat(s)
		if err != nil {
			return nil, err
		}
		res[i] = n
	}
	return res, nil
}

func parseFloat(param any) (float64, error) {
	switch v := param.(type) {
	case string:
		var n float64
		err := parseStrFloat(v, &n)
		return n, err
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, errors.New("invalid parameter")
}

func parseStrFloat(s string, out *float64) error {
	n := 0.0
	sign := 1.0
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		if s[i] == '-' {
			sign = -1
		}
		i++
	}
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		n = n*10 + float64(s[i]-'0')
		i++
	}
	if i < len(s) && s[i] == '.' {
		i++
		place := 0.1
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			n += float64(s[i]-'0') * place
			place *= 0.1
			i++
		}
	}
	*out = n * sign
	if i != len(s) {
		return errors.New("invalid float string: " + s)
	}
	return nil
}

// Turtle state mirrored by the app.
var (
	turtleX, turtleY, turtleAngle float64
	turtleDown                    = true
	turtleColor                   = [3]uint8{20, 20, 20}
	turtleWidth                   = 2.0
)
