package turtle

import (
	"testing"

	"github.com/zrcoder/ndor/cmd/gui/turtle/api"
)

func TestRunCode(t *testing.T) {
	api.CmdQueue = nil
	code := `for i := 0; i < 4; i++ {
	step(5)
	right(90)
}`
	if err := runCode(code); err != nil {
		t.Fatalf("runCode: %v", err)
	}
	if got, want := len(api.CmdQueue), 8; got != want {
		t.Fatalf("got %d commands, want %d", got, want)
	}
}

func TestRunCodeCommandStyle(t *testing.T) {
	api.CmdQueue = nil
	code := `color "red"
step 10
right 90
step -5
left 45`
	if err := runCode(code); err != nil {
		t.Fatalf("runCode: %v", err)
	}
	if got, want := len(api.CmdQueue), 5; got != want {
		t.Fatalf("got %d commands, want %d", got, want)
	}
}

func TestStepSign(t *testing.T) {
	api.CmdQueue = nil
	if err := runCode("step 3\nstep -3"); err != nil {
		t.Fatalf("runCode: %v", err)
	}
	if len(api.CmdQueue) != 2 {
		t.Fatalf("got %d commands, want 2", len(api.CmdQueue))
	}
	if c := api.CmdQueue[0]; c.Kind != "step" || c.Val != 3 {
		t.Fatalf("step 3 => %+v, want step 3", c)
	}
	if c := api.CmdQueue[1]; c.Kind != "step" || c.Val != -3 {
		t.Fatalf("step -3 => %+v, want step -3", c)
	}
}

func TestExamples(t *testing.T) {
	for _, ex := range examples {
		api.CmdQueue = nil
		if err := runCode(ex.code); err != nil {
			t.Errorf("%s: %v", ex.name, err)
			continue
		}
		if len(api.CmdQueue) == 0 {
			t.Errorf("%s: no commands produced", ex.name)
		}
	}
}
