package turtle

import (
	"math"

	api "github.com/zrcoder/ndor/api/turtle"
	"github.com/zrcoder/ndor/internal/turtle/xgo"
	"github.com/zrcoder/ndor/internal/xgoerr"
)

const (
	cell               = 10
	stepPixelsPerFrame = 4
)

// op is one micro-operation replayed by the app.
type op struct {
	kind  string
	val   float32
	color [3]uint8
}

// Program is a compiled turtle script that can be played back step by step.
type Program struct {
	state state
	ops   []op
	idx   int
}

// New returns an empty program with a cleared turtle.
func New() *Program {
	var s state
	s.clear()
	return &Program{state: s}
}

// Compile runs xgo code and queues its drawing operations.
func Compile(code string) (*Program, *xgoerr.LineError) {
	p := New()
	api.CmdQueue = nil
	if err := xgo.Run(code); err != nil {
		return p, err
	}

	for _, cmd := range api.CmdQueue {
		switch cmd.Kind {
		case "step":
			// Split each step into small moves so the animation is smooth.
			dist := float32(cmd.Val) * cell
			n := max(int(math.Ceil(float64(dist)/float64(stepPixelsPerFrame))), 1)
			step := dist / float32(n)
			for range n {
				p.ops = append(p.ops, op{kind: "step", val: step})
			}
		case "turn":
			p.ops = append(p.ops, op{kind: "turn", val: float32(cmd.Val)})
		case "penUp":
			p.ops = append(p.ops, op{kind: "penUp"})
		case "penDown":
			p.ops = append(p.ops, op{kind: "penDown"})
		case "setColor":
			p.ops = append(p.ops, op{kind: "setColor", color: cmd.Color})
		case "setWidth":
			p.ops = append(p.ops, op{kind: "setWidth", val: float32(cmd.Val)})
		case "clear":
			p.ops = append(p.ops, op{kind: "clear"})
		case "home":
			p.ops = append(p.ops, op{kind: "home"})
		}
	}
	return p, nil
}

// Done reports whether all queued operations have been applied.
func (p *Program) Done() bool { return p.idx >= len(p.ops) }

// Advance applies the next queued operation.
func (p *Program) Advance() {
	if p.Done() {
		return
	}
	o := p.ops[p.idx]
	p.idx++

	switch o.kind {
	case "step":
		p.state.step(o.val)
	case "turn":
		p.state.turn(o.val)
	case "penUp":
		p.state.penUp()
	case "penDown":
		p.state.penDown()
	case "setColor":
		p.state.setColor(o.color[0], o.color[1], o.color[2])
	case "setWidth":
		p.state.setWidth(o.val)
	case "clear":
		p.state.clear()
	case "home":
		p.state.home()
	}
}
