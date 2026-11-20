package comp

import (
	"time"

	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// PaintingDuration is how long the painting overlay shows.
const PaintingDuration = 1200 * time.Millisecond

// PaintingActive reports whether the painting overlay shows until.
func PaintingActive(c *ui.Context, until time.Time) bool {
	return !c.Preferences().ReduceMotion && time.Now().Before(until)
}

// Painting draws the painting overlay across the window.
func Painting(c *ui.Context, left, right *ui.Bitmap) {
	t := c.Theme()
	ui.Row(c.Key("painting")).Absolute().Top(0).Bottom(0).Left(0).Right(0).Center().PassThrough().Children(func() {
		ui.Row(c).Radius(t.Space(2.5)).
			Background(ui.RGBA(255, 255, 255, 0.6)).Border(1, t.Border).
			Shadow(0, 10, 30, 0, ui.RGBA(0, 0, 0, 0.3)).Label("Painting").Children(func() {
			if left != nil {
				ui.Image(c, left).MaxHeightPercent(31)
			}

			loading.Blocks(c, loading.Options{
				Size:  80,
				Color: t.Accent,
			}).Key("painting-spinner")

			if right != nil {
				ui.Image(c, right).MaxHeightPercent(28)
			}
		})
	})
}
