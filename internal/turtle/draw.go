package turtle

import "github.com/egoist/mygo/ui"

// Draw renders the grid, trail and the gopher sprite centered in r.
func (p *Program) Draw(painter *ui.Painter, r ui.Rect, gopher *ui.Bitmap) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2

	// Grid lines darken with distance from the origin; the axes are darkest.
	gridLight := ui.RGB(240, 240, 240)
	gridMid := ui.RGB(220, 220, 220)
	gridCenter := ui.RGB(150, 150, 150)

	strokeGrid := func(x0, y0, x1, y1 float32, col ui.Color) {
		ln := ui.Path{}
		ln.MoveTo(x0, y0).LineTo(x1, y1)
		painter.StrokePath(&ln, 1, col)
	}

	nx := int(cx / cell)
	for k := -nx; k <= nx; k++ {
		x := cx + float32(k)*cell
		if x < r.X || x > r.X+r.W {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(x, r.Y, x, r.Y+r.H, col)
	}

	ny := int(cy / cell)
	for k := -ny; k <= ny; k++ {
		y := cy + float32(k)*cell
		if y < r.Y || y > r.Y+r.H {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(r.X, y, r.X+r.W, y, col)
	}

	sx := func(q pt) float32 { return cx + q.X }
	sy := func(q pt) float32 { return cy - q.Y }

	strokeSeg := func(s seg) {
		if len(s.pts) < 2 {
			return
		}
		ln := ui.Path{}
		ln.MoveTo(sx(s.pts[0]), sy(s.pts[0]))
		for _, q := range s.pts[1:] {
			ln.LineTo(sx(q), sy(q))
		}
		painter.StrokePath(&ln, s.width, ui.RGB(s.color[0], s.color[1], s.color[2]))
	}
	for _, s := range p.state.segments {
		strokeSeg(s)
	}
	if p.state.cur != nil {
		strokeSeg(*p.state.cur)
	}

	px, py := sx(pt{p.state.x, p.state.y}), sy(pt{p.state.x, p.state.y})
	drawGopher(painter, gopher, px, py, p.state.angle)
}

const gopherDrawSize = 30

// drawGopher draws the gopher sprite at (x, y) rotated by angle.
func drawGopher(p *ui.Painter, sprite *ui.Bitmap, x, y, angle float32) {
	// ui.Rotation is clockwise in screen space; turtle angle is math space (+Y up, 0 right).
	// The gopher faces up by default, so screen rotation is 90 - angle.
	screenAngle := 90 - angle

	p.Transform(ui.Rotation(screenAngle), x, y, func() {
		dst := ui.Rect{X: x - gopherDrawSize/2, Y: y - gopherDrawSize/2, W: gopherDrawSize, H: gopherDrawSize}
		p.ImageRegion(sprite, ui.Rect{}, dst, ui.ImageOptions{})
	})
}
