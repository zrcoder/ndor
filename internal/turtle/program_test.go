package turtle

import (
	"math"
	"testing"
)

func TestColorAlpha(t *testing.T) {
	p, lineErr := Compile("color 255, 0, 0, 128\n")
	if lineErr != nil {
		t.Fatal(lineErr)
	}
	for !p.Done() {
		p.Advance()
	}
	if want := [4]uint8{255, 0, 0, 128}; p.state.penColor != want {
		t.Fatalf("pen color = %v, want %v", p.state.penColor, want)
	}

	p, lineErr = Compile("color \"red\"\n")
	if lineErr != nil {
		t.Fatal(lineErr)
	}
	for !p.Done() {
		p.Advance()
	}
	if p.state.penColor[3] != 255 {
		t.Fatalf("named color alpha = %d, want 255", p.state.penColor[3])
	}
}

func TestCompileSubSteps(t *testing.T) {
	p, lineErr := Compile("step 1000000\n")
	if lineErr != nil {
		t.Fatal(lineErr)
	}
	if len(p.ops) != maxSubOps {
		t.Fatalf("sub-op count = %d, want %d", len(p.ops), maxSubOps)
	}
	var total float32
	for _, op := range p.ops {
		total += op.val
	}
	want := float32(1000000) * cell
	if math.Abs(float64(total-want)) > float64(want)*1e-4 {
		t.Fatalf("step total = %v, want %v", total, want)
	}

	p, lineErr = Compile("step -10\n")
	if lineErr != nil {
		t.Fatal(lineErr)
	}
	if len(p.ops) != 25 || p.ops[0].val != -4 {
		t.Fatalf("negative sub-steps = %+v, want 25 steps of -4px", p.ops[:min(len(p.ops), 3)])
	}
}
