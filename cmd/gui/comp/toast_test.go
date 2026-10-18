package comp

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestToast(t *testing.T) {
	dismissed := false
	tt := ui.NewTester(func(c *ui.Context) {
		Toast(c, "Hello World", func() {
			dismissed = true
		})
	}, 400, 300)

	if !tt.HasText("Hello World") {
		t.Errorf("expected toast message, got: %q", tt.Texts())
	}

	// Click to dismiss
	tt.ClickAt(200, 150)
	if !dismissed {
		t.Error("expected toast to be dismissed on click")
	}
}

func TestToastWithTitle(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ToastWithTitle(c, "Title", "Message", nil)
	}, 400, 300)

	if !tt.HasText("Title: Message") {
		t.Errorf("expected 'Title: Message', got: %q", tt.Texts())
	}
}

func TestToastErrorAtLine(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ToastErrorAtLine(c, 5, "undefined: foo", nil)
	}, 400, 300)

	if !tt.HasText("Error at line 5: undefined: foo") {
		t.Errorf("expected error at line, got: %q", tt.Texts())
	}
}
