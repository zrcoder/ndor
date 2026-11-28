package turtle

import (
	"math"
	"testing"
)

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
