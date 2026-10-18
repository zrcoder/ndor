package component

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var testCopyIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.3"><rect x="5.5" y="5.5" width="8" height="8" rx="1.6"/><path d="M10.5 3.6V3.2A1.7 1.7 0 0 0 8.8 1.5H4.2A1.7 1.7 0 0 0 2.5 3.2v4.6a1.7 1.7 0 0 0 1.7 1.7h.4"/></svg>`))

var testCheckIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M3.5 8.5l3 3 6-7"/></svg>`))

func TestIconButton(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		IconButton(c, testCopyIcon, nil, false, func() {
			clicks++
		})
	}, 100, 50)

	// Click at button location (center)
	tt.ClickAt(50, 25)
	if clicks != 1 {
		t.Errorf("expected 1 click, got %d", clicks)
	}
}

func TestIconButtonWithCheck(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		IconButton(c, testCopyIcon, testCheckIcon, true, func() {
			clicks++
		})
	}, 100, 50)

	tt.ClickAt(50, 25)
	if clicks != 1 {
		t.Errorf("expected 1 click, got %d", clicks)
	}
}
