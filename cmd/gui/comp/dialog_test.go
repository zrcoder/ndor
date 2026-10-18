package comp

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestShowDialogWithImage(t *testing.T) {
	open := true
	clickedBtn := -1
	opt := DialogOption{
		Open:    &open,
		Buttons: []string{"Left", "Right"},
		OnClick: func(btnIndex int) {
			clickedBtn = btnIndex
		},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		ShowDialogWithImage(c, opt)
	}, 400, 300)

	// Verify buttons are shown
	if !tt.HasText("Left") || !tt.HasText("Right") {
		t.Errorf("expected buttons, got: %q", tt.Texts())
	}

	// Click the left button
	if err := tt.Click("Left"); err != nil {
		t.Fatal(err)
	}
	if clickedBtn != 0 {
		t.Errorf("expected clickedBtn 0, got %d", clickedBtn)
	}
}

func TestShowDialogWithImageDismissed(t *testing.T) {
	open := true
	outsideClicked := false

	opt := DialogOption{
		Open:    &open,
		Buttons: []string{"Button"},
		OnClick: func(btnIndex int) {
			outsideClicked = true
		},
	}
	_ = ui.NewTester(func(c *ui.Context) {
		ShowDialogWithImage(c, opt)
	}, 400, 300)

	// Initial state: outsideClicked should be false since we haven't clicked outside yet
	if outsideClicked {
		t.Error("outside click should not trigger before modal is dismissed")
	}
}
