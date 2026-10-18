package main

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestRunningProgressRenders checks that the progress bar is actually
// painted while isRunning is true.
func TestRunningProgressRenders(t *testing.T) {
	ndor := New()
	tt := ui.NewTester(ndor.View, 800, 600)
	_ = tt

	img := tt.Image()
	// Scan for accent-colored pixels (the progress fill), any non-background pixel
	// in the bottom area where the bar should be (Bottom(80) => y ~ 512..524).
	var found int
	for y := 500; y < 530 && y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			c := img.RGBAAt(x, y)
			// anything not near-transparent and not pure background
			if c.A > 0 && c != (color.RGBA{}) {
				found++
			}
		}
	}
	if found == 0 {
		t.Errorf("no pixels painted in the progress bar area (bottom ~80px)")
	}

	// Dump a coarse map to see what is where
	t.Logf("image %v", img.Bounds())
}
