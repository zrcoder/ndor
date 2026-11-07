package internal

import "testing"

func TestTurtleForward(t *testing.T) {
	ResetTurtle(400, 400)
	Turtle.Forward(100)
	if got := len(TurtleSteps); got != 1 {
		t.Fatalf("expected 1 step, got %d", got)
	}
	s := TurtleSteps[0]
	if s.X1 != 200 || s.Y1 != 200 || s.X2 != 200 || s.Y2 != 100 {
		t.Fatalf("unexpected step: %+v", s)
	}

	Turtle.Left(90)
	Turtle.Forward(100)
	if Turtle.heading != -90 {
		t.Fatalf("unexpected heading: %v", Turtle.heading)
	}

	Turtle.Pen(false)
	Turtle.Forward(100)
	if len(TurtleSteps) != 2 {
		t.Fatalf("pen up should not add a step, got %d", len(TurtleSteps))
	}
	x, y, _ := Turtle.End()
	if x != 0 || y != 100 {
		t.Fatalf("unexpected pen-up end: %v, %v", x, y)
	}
}

func TestTurtleEnd(t *testing.T) {
	ResetTurtle(400, 400)
	Turtle.Jump(10, 20)
	x, y, h := Turtle.End()
	if x != 210 || y != 180 || h != 0 {
		t.Fatalf("unexpected end: %v, %v, %v", x, y, h)
	}
}
