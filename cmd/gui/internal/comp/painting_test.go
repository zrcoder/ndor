package comp

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestPaintingDrawsCentered(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Painting(c, nil, nil)
	}, 200, 200)
	tt.Frame()

	r, ok := tt.Find("Painting")
	if !ok {
		t.Fatal("painting does not show")
	}
	if cx, cy := r.X+r.W/2, r.Y+r.H/2; cx != 100 || cy != 100 {
		t.Fatalf("painting is not centered in the window: center (%v, %v)", cx, cy)
	}
}

func TestPaintingActive(t *testing.T) {
	until := time.Now().Add(time.Hour)
	active := false
	tt := ui.NewTester(func(c *ui.Context) {
		active = PaintingActive(c, until)
	}, 100, 100)

	tt.Frame()
	if !active {
		t.Fatal("painting not active before its deadline")
	}

	until = time.Now().Add(-time.Hour)
	tt.Frame()
	if active {
		t.Fatal("painting active after its deadline")
	}

	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	until = time.Now().Add(time.Hour)
	tt.Frame()
	if active {
		t.Fatal("painting active under reduced motion")
	}
}
