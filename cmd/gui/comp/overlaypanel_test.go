package comp

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var testCloseIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M4 4l8 8M12 4l-8 8"/></svg>`))

func TestOverlayPanel(t *testing.T) {
	selected := ""
	tt := ui.NewTester(func(c *ui.Context) {
		cfg := OverlayPanelOption{
			Title:     "Test Panel",
			Width:     180,
			CloseIcon: testCheckIcon,
			ItemFn: func(item string) {
				selected = item
			},
		}
		OverlayPanel(c, cfg, []string{"Item 1", "Item 2"})
	}, 400, 300)

	// Verify panel title is shown
	if !tt.HasText("Test Panel") {
		t.Errorf("expected panel title, got: %q", tt.Texts())
	}

	// Verify items are shown
	if !tt.HasText("Item 1") || !tt.HasText("Item 2") {
		t.Errorf("expected items, got: %q", tt.Texts())
	}

	// Click on an item
	if err := tt.Click("Item 1"); err != nil {
		t.Fatal(err)
	}
	if selected != "Item 1" {
		t.Errorf("expected 'Item 1', got %q", selected)
	}
}

func TestOverlayPanelCloseOnClickOutside(t *testing.T) {
	closed := false
	tt := ui.NewTester(func(c *ui.Context) {
		opt := OverlayPanelOption{
			Title:          "Test",
			Width:          180,
			OnClose:        func() { closed = true },
			OnClickOutside: func() { closed = true },
			CloseIcon:      testCloseIcon,
			ItemFn:         func(item string) {},
		}
		OverlayPanel(c, opt, []string{"Item"})
	}, 400, 300)

	// Click outside the panel (at top-left corner where there's no panel)
	tt.ClickAt(5, 5)
	if !closed {
		t.Error("expected panel to close on click outside")
	}
}
