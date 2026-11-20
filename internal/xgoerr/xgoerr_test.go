package xgoerr

import "testing"

func TestFindLine(t *testing.T) {
	code := `line 1
line 2
...
color "unkown"
...`
	if got := FindLine(code, "color"); got != 4 {
		t.Errorf("FindLine() = %v, want 4", got)
	}
	if got := FindLine(code, "circle"); got != -1 {
		t.Errorf("FindLine() = %v, want -1", got)
	}
}
