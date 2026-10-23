package comp

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestHelpShowsAfterOpen(t *testing.T) {
	h := NewHelp(HelpConfig{
		ExampleNames: []string{"one"},
		Examples:     map[string]string{"one": "code"},
		LoadDoc:      func(bool, int) string { return "doc" },
	})

	tt := ui.NewTester(func(c *ui.Context) { h.View(c) }, 400, 300)
	if tt.HasText("Examples") {
		t.Fatal("help shows before Open")
	}

	h.Open()
	tt.Frame()
	if !tt.HasText("Examples") || !tt.HasText("Document") {
		t.Fatal("help dialog does not show after Open")
	}
}
