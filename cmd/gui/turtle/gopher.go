package turtle

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

//go:embed gopher_walk.svg
var gopherFS embed.FS

// gopher_walk.svg is trimmed to its first frame (50x50).
const gopherDrawSize = 30

var gopherWalk *ui.Bitmap

func initGopher() error {
	walkData, err := gopherFS.ReadFile("gopher_walk.svg")
	if err != nil {
		return err
	}

	walkSVG, err := ui.ParseSVG(walkData)
	if err != nil {
		return err
	}

	// Rasterize once; the bitmap is reused.
	gopherWalk = ui.NewBitmap(walkSVG.Image(50, 50, nil))
	return nil
}

// drawGopher draws the static gopher at (x, y) rotated by angle.
func drawGopher(p *ui.Painter, x, y float32, angle float32) {
	if gopherWalk == nil {
		return
	}

	// ui.Rotation is clockwise in screen space; turtle angle is math space (+Y up, 0 right).
	// The gopher faces up by default, so screen rotation is 90 - angle.
	screenAngle := 90 - angle

	p.Transform(ui.Rotation(screenAngle), x, y, func() {
		dst := ui.Rect{X: x - gopherDrawSize/2, Y: y - gopherDrawSize/2, W: gopherDrawSize, H: gopherDrawSize}
		p.ImageRegion(gopherWalk, ui.Rect{}, dst, ui.ImageOptions{})
	})
}
