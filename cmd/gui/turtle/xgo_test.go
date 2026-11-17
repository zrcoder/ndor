package turtle

import (
	"testing"

	"github.com/zrcoder/ndor/cmd/gui/turtle/api"
)

func TestRunCode(t *testing.T) {
	api.CmdQueue = nil
	code := `for i := 0; i < 4; i++ {
	forward(5)
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
forward 10
right 90
backward 5
left 45`
	if err := runCode(code); err != nil {
		t.Fatalf("runCode: %v", err)
	}
	if got, want := len(api.CmdQueue), 5; got != want {
		t.Fatalf("got %d commands, want %d", got, want)
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
