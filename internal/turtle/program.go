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
	turnDegPerFrame    = 14
	// maxSubOps bounds compile-time allocation for huge step/turn values.
	maxSubOps = 1000
)

// op is one micro-operation replayed by the app.
type op struct {
	kind  string
	val   float32
	color [4]uint8
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
		val, perFrame := float32(cmd.Val), float64(stepPixelsPerFrame)
		switch cmd.Kind {
		case "step":
			val *= cell
		case "turn":
			perFrame = turnDegPerFrame
		default:
			p.ops = append(p.ops, op{kind: cmd.Kind, val: float32(cmd.Val), color: cmd.Color})
			continue
		}
		n := max(int(math.Ceil(math.Abs(float64(val))/perFrame)), 1)
		n = min(n, maxSubOps)
		for range n {
			p.ops = append(p.ops, op{kind: cmd.Kind, val: val / float32(n)})
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
		p.state.setColor(o.color)
	case "setWidth":
		p.state.setWidth(o.val)
	case "clear":
		p.state.clear()
	case "home":
		p.state.home()
	}
}
