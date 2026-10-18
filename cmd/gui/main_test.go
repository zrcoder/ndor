package main

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestView(t *testing.T) {
	ndor := New()
	ui.NewTester(ndor.View, 800, 600)
	// Test that the view renders without crashing
	// Buttons use labeled images, so we verify the dialog image exists
	if ndor.helpDialogImg == nil {
		t.Errorf("expected help dialog image to be loaded")
	}
}
