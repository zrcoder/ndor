package comp

import (
	"time"

	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
)

// PaintingDuration is how long the painting overlay shows.
const PaintingDuration = 1200 * time.Millisecond

// paintingToast names the toast that shows the painting overlay, so the
// viewport can tell it apart from other toasts.
const paintingToast = "painting"

// ShowPainting shows the painting overlay. It is a toast, so its own
// lifetime decides when the overlay goes away.
func ShowPainting(c *ui.Context) {
	if c.Preferences().ReduceMotion {
		return
	}
	c.AddToast(ui.Toast{
		ID:      paintingToast,
		Title:   "Painting",
		Timeout: PaintingDuration,
	})
}

// PaintingViewport draws the painting toast centered over the window, and
// reports whether it shows, so the caller can keep the result hidden until
// the painting is over. Call it once in the view, where the window's toasts
// show.
func PaintingViewport(c *ui.Context, left, right *ui.Bitmap) bool {
	showing := false
	ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.AlignItems(ui.Center).Gap(8).Padding(16)
		for _, t := range toasts {
			if t.ID == paintingToast {
				showing = true
				painting(c, left, right)
				continue
			}
			// Other toasts keep a plain look at the bottom.
			toast := ui.ToastBase(c, t)
			toast.Root.Row().Gap(12).Padding(10, 14).Radius(8).
				Background(c.Theme().Surface).Children(func() {
				if t.Title != "" {
					ui.Text(c, t.Title)
				}
				if t.Description != "" {
					ui.Text(c, t.Description)
				}
			})
		}
	})
	return showing
}

func painting(c *ui.Context, left, right *ui.Bitmap) {
	ui.Row(c).Key("painting").Absolute().Top(0).Bottom(0).Left(0).Right(0).
		Center().Gap(6).PassThrough().Children(func() {
		if left != nil {
			ui.Image(c, left).MaxHeightPercent(11)
		}

		loading.Blocks(c, loading.Options{
			Size:  80,
			Color: c.Theme().Accent,
		}).Label("Painting")

		if right != nil {
			ui.Image(c, right).MaxHeightPercent(8)
		}
	})
}
