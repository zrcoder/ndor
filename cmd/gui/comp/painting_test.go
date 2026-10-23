package comp

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestPaintingShowsAsToastAndCloses(t *testing.T) {
	show := false
	showing := false
	tt := ui.NewTester(func(c *ui.Context) {
		if show {
			show = false
			ShowPainting(c)
		}
		showing = PaintingViewport(c, nil, nil)
	}, 200, 200)

	if showing {
		t.Fatal("painting reports showing before ShowPainting")
	}
	if _, ok := tt.Find("Painting"); ok {
		t.Fatal("painting shows before ShowPainting")
	}

	show = true
	tt.Frame()
	r, ok := tt.Find("Painting")
	if !ok || !showing {
		t.Fatal("painting does not show the spinner after ShowPainting")
	}
	if cx, cy := r.X+r.W/2, r.Y+r.H/2; cx != 100 || cy != 100 {
		t.Fatalf("painting is not centered in the window: center (%v, %v)", cx, cy)
	}

	tt.Advance(PaintingDuration + 100*time.Millisecond)
	if showing {
		t.Fatal("painting still reports showing after its toast closes")
	}
	if _, ok := tt.Find("Painting"); ok {
		t.Fatal("painting spinner still shows after its toast closes")
	}

	// Showing it again starts a fresh toast.
	show = true
	tt.Frame()
	if !showing {
		t.Fatal("painting does not report showing again")
	}
	if _, ok := tt.Find("Painting"); !ok {
		t.Fatal("painting does not show again")
	}
}

func TestPaintingSkipsReducedMotion(t *testing.T) {
	show := false
	showing := false
	tt := ui.NewTester(func(c *ui.Context) {
		if show {
			show = false
			ShowPainting(c)
		}
		showing = PaintingViewport(c, nil, nil)
	}, 200, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	show = true
	tt.Frame()

	if showing {
		t.Fatal("painting reports showing under reduced motion")
	}
	if _, ok := tt.Find("Painting"); ok {
		t.Fatal("painting shows under reduced motion")
	}
}
