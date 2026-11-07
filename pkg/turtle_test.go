package pkg

import (
	"testing"
)

func TestTurtleGen(t *testing.T) {
	plan, err := TurtleGen(100, 100, `context 100, 100
forward 50`)
	if err != nil {
		t.Fatal(err)
	}
	if plan.W != 100 || plan.H != 100 {
		t.Fatalf("unexpected size: %dx%d", plan.W, plan.H)
	}
	if len(plan.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(plan.Steps))
	}
	s := plan.Steps[0]
	if s.X1 != 50 || s.Y1 != 50 || s.X2 != 50 || s.Y2 != 0 {
		t.Fatalf("unexpected step: %+v", s)
	}
	if s.A != 0xff {
		t.Fatalf("expected opaque color, got %+v", s)
	}

	_, err = TurtleGen(100, 100, `context 100, 100
forward`)
	if err == nil {
		t.Fatal("expected an error for missing argument")
	}

	plan, err = TurtleGen(100, 100, `context 100, 100
forward 10`)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Steps[0].R != 0 || plan.Steps[0].G != 0 || plan.Steps[0].B != 0 {
		t.Fatalf("expected black pen, got %+v", plan.Steps[0])
	}
}
