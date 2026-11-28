package ndor

import (
	"fmt"
	"image"
	"math"

	"github.com/egoist/mygo/ui"
)

// Radians converts degrees to radians.
func Radians(angle float64) float64 {
	return angle * math.Pi / 180
}

// matrix is a 2D affine transform laid out like gg's Matrix.
type matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
}

func (a matrix) multiply(b matrix) matrix {
	return matrix{
		XX: a.XX*b.XX + a.YX*b.XY,
		YX: a.XX*b.YX + a.YX*b.YY,
		XY: a.XY*b.XX + a.YY*b.XY,
		YY: a.XY*b.YX + a.YY*b.YY,
		X0: a.X0*b.XX + a.Y0*b.XY + b.X0,
		Y0: a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

func (a matrix) transformPoint(x, y float64) (tx, ty float64) {
	return a.XX*x + a.XY*y + a.X0, a.YX*x + a.YY*y + a.Y0
}

// state is what Push saves and Pop restores (not the current path).
type state struct {
	matrix      matrix
	fillColor   ui.Color
	strokeColor ui.Color
	textColor   ui.Color
	dashes      []float64
	lineWidth   float64
}

// drawOp is one recorded command, replayed in order by Image.
type drawOp interface {
	draw(p *ui.Painter)
}

type fillOp struct {
	path  ui.Path
	color ui.Color
}

func (o fillOp) draw(p *ui.Painter) { p.FillPath(&o.path, o.color) }

type strokeOp struct {
	path   ui.Path
	color  ui.Color
	width  float32
	dashes []float32
}

func (o strokeOp) draw(p *ui.Painter) {
	o.path.Dash(0, o.dashes...)
	p.StrokePath(&o.path, o.width, o.color)
}

type clipOp struct {
	path ui.Path
}

func (o clipOp) draw(p *ui.Painter) {}

type textOp struct {
	text  string
	x, y  float64
	color ui.Color
	scale float64
	angle float64
}

func (o textOp) size() float32 {
	size := float32(13 * o.scale)
	if size < 1 {
		size = 13
	}
	return size
}

func (o textOp) draw(p *ui.Painter) {
	size := o.size()
	x, y := float32(o.x), float32(o.y)
	top := y - size*0.8
	if o.angle != 0 {
		p.Transform(ui.Rotation(float32(o.angle*180/math.Pi)), x, y, func() {
			p.Text(x, top, o.text, size, o.color)
		})
		return
	}
	p.Text(x, top, o.text, size, o.color)
}

// bounds estimates the box the text covers, rotation included.
func (o textOp) bounds() ui.Rect {
	size := float64(o.size())
	w := 0.0
	for _, r := range o.text {
		if r < 0x2E80 {
			w += size * 0.55
		} else {
			w += size
		}
	}
	parts := [4][2]float64{{0, -size}, {w, -size}, {w, size * 0.35}, {0, size * 0.35}}
	cos, sin := math.Cos(o.angle), math.Sin(o.angle)
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, q := range parts {
		x, y := o.x+q[0]*cos-q[1]*sin, o.y+q[0]*sin+q[1]*cos
		x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x), max(y1, y)
	}
	return ui.Rect{X: float32(x0), Y: float32(y0), W: float32(x1 - x0), H: float32(y1 - y0)}
}

type clearOp struct {
	color ui.Color
}

func (o clearOp) draw(p *ui.Painter) {}

// Context records the ndor API and renders it with mygo's painter.
type Context struct {
	matrix      matrix
	fillColor   ui.Color
	strokeColor ui.Color
	textColor   ui.Color
	dashes      []float64
	lineWidth   float64
	stack       []state

	path       ui.Path
	hasCurrent bool
	startX     float64
	startY     float64
	curX       float64
	curY       float64

	pathBox box // bounds of the path being built
	clipBox box // bounds of the active clip, if any
	ink     box // bounds of everything drawn

	ops []drawOp
}

// box tracks the four edges of a region.
type box struct {
	x0, y0, x1, y1 float32
	ok             bool
}

func (b *box) add(x, y float32) {
	if !b.ok {
		b.x0, b.y0, b.x1, b.y1, b.ok = x, y, x, y, true
		return
	}
	b.x0, b.y0, b.x1, b.y1 = min(b.x0, x), min(b.y0, y), max(b.x1, x), max(b.y1, y)
}

func (b *box) addRect(r ui.Rect) {
	b.add(r.X, r.Y)
	b.add(r.X+r.W, r.Y+r.H)
}

// addCurve widens b to the bezier through the x,y pairs, sampled because
// its control points may lie outside the curve.
func (b *box) addCurve(pts ...float64) {
	const n = 32
	for i := range n + 1 {
		t := float64(i) / n
		u := 1 - t
		q := append([]float64(nil), pts...)
		for len(q) > 2 {
			for j := range len(q) - 2 {
				q[j] = u*q[j] + t*q[j+2]
			}
			q = q[:len(q)-2]
		}
		b.add(float32(q[0]), float32(q[1]))
	}
}

// NewContext returns a context whose origin is the picture center, +Y up.
func NewContext() *Context {
	return &Context{
		matrix:      matrix{XX: 1, YY: -1},
		fillColor:   ui.RGB(255, 255, 255),
		strokeColor: ui.RGB(0, 0, 0),
		textColor:   ui.Transparent,
		lineWidth:   1,
	}
}

// SetRGBA255 sets the color for subsequent fills, strokes and text.
func (c *Context) SetRGBA255(r, g, b, a int) {
	rgba := ui.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	c.fillColor, c.strokeColor, c.textColor = rgba, rgba, rgba
}

func (c *Context) newSubPath() {
	c.hasCurrent = false
}

func (c *Context) MoveTo(x, y float64) {
	x, y = c.matrix.transformPoint(x, y)
	c.path.MoveTo(float32(x), float32(y))
	c.startX, c.startY, c.curX, c.curY = x, y, x, y
	c.pathBox.add(float32(x), float32(y))
	c.hasCurrent = true
}

func (c *Context) LineTo(x, y float64) {
	if !c.hasCurrent {
		c.MoveTo(x, y)
		return
	}
	x, y = c.matrix.transformPoint(x, y)
	c.path.LineTo(float32(x), float32(y))
	c.curX, c.curY = x, y
	c.pathBox.add(float32(x), float32(y))
}

func (c *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !c.hasCurrent {
		c.MoveTo(x1, y1)
	}
	x1, y1 = c.matrix.transformPoint(x1, y1)
	x2, y2 = c.matrix.transformPoint(x2, y2)
	c.path.QuadTo(float32(x1), float32(y1), float32(x2), float32(y2))
	c.pathBox.addCurve(c.curX, c.curY, x1, y1, x2, y2)
	c.curX, c.curY = x2, y2
}

func (c *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !c.hasCurrent {
		c.MoveTo(x1, y1)
	}
	x1, y1 = c.matrix.transformPoint(x1, y1)
	x2, y2 = c.matrix.transformPoint(x2, y2)
	x3, y3 = c.matrix.transformPoint(x3, y3)
	c.path.CubeTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(x3), float32(y3))
	c.pathBox.addCurve(c.curX, c.curY, x1, y1, x2, y2, x3, y3)
	c.curX, c.curY = x3, y3
}

func (c *Context) closePath() {
	if c.hasCurrent {
		c.path.LineTo(float32(c.startX), float32(c.startY))
		c.path.Close()
		c.curX, c.curY = c.startX, c.startY
	}
}

func (c *Context) clearPath() {
	c.path = ui.Path{}
	c.pathBox = box{}
	c.hasCurrent = false
}

func (c *Context) DrawRectangle(x, y, w, h float64) {
	c.newSubPath()
	c.MoveTo(x, y)
	c.LineTo(x+w, y)
	c.LineTo(x+w, y+h)
	c.LineTo(x, y+h)
	c.closePath()
}

func (c *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	c.newSubPath()
	c.MoveTo(x1, y0)
	c.LineTo(x2, y0)
	c.DrawArc(x2, y1, r, Radians(270), Radians(360))
	c.LineTo(x3, y2)
	c.DrawArc(x2, y2, r, Radians(0), Radians(90))
	c.LineTo(x1, y3)
	c.DrawArc(x1, y2, r, Radians(90), Radians(180))
	c.LineTo(x0, y1)
	c.DrawArc(x1, y1, r, Radians(180), Radians(270))
	c.closePath()
}

func (c *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := range n {
		p1 := float64(i) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos((a1+a2)/2)
		y1 := y + ry*math.Sin((a1+a2)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if c.hasCurrent {
				c.LineTo(x0, y0)
			} else {
				c.MoveTo(x0, y0)
			}
		}
		c.QuadraticTo(cx, cy, x2, y2)
	}
}

func (c *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	c.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

func (c *Context) DrawEllipse(x, y, rx, ry float64) {
	c.newSubPath()
	c.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	c.closePath()
}

func (c *Context) DrawCircle(x, y, r float64) {
	c.newSubPath()
	c.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
	c.closePath()
}

func (c *Context) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	if n%2 == 0 {
		rotation += angle / 2
	}
	c.newSubPath()
	for i := range n {
		a := rotation + angle*float64(i)
		c.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	c.closePath()
}

func (c *Context) SetDash(dashes ...float64) {
	c.dashes = dashes
}

func (c *Context) SetLineWidth(width float64) {
	c.lineWidth = width
}

func (c *Context) DrawString(s string, x, y float64) {
	tx, ty := c.matrix.transformPoint(x, y)
	angle := math.Atan2(c.matrix.YX, c.matrix.XX)
	scale := math.Hypot(c.matrix.XX, c.matrix.YX)
	op := textOp{text: s, x: tx, y: ty, color: c.textColor, scale: scale, angle: angle}
	c.ops = append(c.ops, op)
	c.ink.addRect(op.bounds())
}

// pathRect returns the path's box widened by pad, cut to the active clip.
func (c *Context) pathRect(pad float32) (ui.Rect, bool) {
	b := c.pathBox
	if !b.ok {
		return ui.Rect{}, false
	}
	r := ui.Rect{X: b.x0 - pad, Y: b.y0 - pad, W: b.x1 - b.x0 + 2*pad, H: b.y1 - b.y0 + 2*pad}
	if c.clipBox.ok {
		x0, y0 := max(r.X, c.clipBox.x0), max(r.Y, c.clipBox.y0)
		x1, y1 := min(r.X+r.W, c.clipBox.x1), min(r.Y+r.H, c.clipBox.y1)
		if x1 <= x0 || y1 <= y0 {
			return ui.Rect{}, false
		}
		r = ui.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
	}
	return r, true
}

func (c *Context) Fill() {
	c.FillPreserve()
	c.clearPath()
}

func (c *Context) FillPreserve() {
	if r, ok := c.pathRect(0); ok {
		c.ink.addRect(r)
	}
	c.ops = append(c.ops, fillOp{path: c.path, color: c.fillColor})
}

func (c *Context) Stroke() {
	c.StrokePreserve()
	c.clearPath()
}

func (c *Context) StrokePreserve() {
	dashes := make([]float32, len(c.dashes))
	for i, v := range c.dashes {
		dashes[i] = float32(v)
	}
	if r, ok := c.pathRect(float32(c.lineWidth / 2)); ok {
		c.ink.addRect(r)
	}
	c.ops = append(c.ops, strokeOp{
		path:   c.path,
		color:  c.strokeColor,
		width:  float32(c.lineWidth),
		dashes: dashes,
	})
}

func (c *Context) Clip() {
	c.ops = append(c.ops, clipOp{path: c.path})
	if c.pathBox.ok {
		if c.clipBox.ok {
			c.clipBox.x0, c.clipBox.y0 = max(c.clipBox.x0, c.pathBox.x0), max(c.clipBox.y0, c.pathBox.y0)
			c.clipBox.x1, c.clipBox.y1 = min(c.clipBox.x1, c.pathBox.x1), min(c.clipBox.y1, c.pathBox.y1)
		} else {
			c.clipBox = c.pathBox
		}
	}
	c.clearPath()
}

func (c *Context) Clear() {
	c.ops = append(c.ops, clearOp{color: c.textColor})
	c.ink = box{}
}

func (c *Context) Push() {
	c.stack = append(c.stack, state{
		matrix:      c.matrix,
		fillColor:   c.fillColor,
		strokeColor: c.strokeColor,
		textColor:   c.textColor,
		dashes:      c.dashes,
		lineWidth:   c.lineWidth,
	})
}

func (c *Context) Pop() {
	if len(c.stack) == 0 {
		return
	}
	s := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	c.matrix = s.matrix
	c.fillColor = s.fillColor
	c.strokeColor = s.strokeColor
	c.textColor = s.textColor
	c.dashes = s.dashes
	c.lineWidth = s.lineWidth
}

func (c *Context) Translate(x, y float64) {
	c.matrix = matrix{XX: 1, YY: 1, X0: x, Y0: y}.multiply(c.matrix)
}

func (c *Context) Scale(x, y float64) {
	c.matrix = matrix{XX: x, YY: y}.multiply(c.matrix)
}

func (c *Context) RotateAbout(angle, x, y float64) {
	c.Translate(x, y)
	cos, sin := math.Cos(angle), math.Sin(angle)
	c.matrix = matrix{XX: cos, YX: sin, XY: -sin, YY: cos}.multiply(c.matrix)
	c.Translate(-x, -y)
}

// Image renders the recorded drawing at the exact size of its content.
// It refuses drawings over maxSide a side, where the canvas would take GBs.
func (c *Context) Image() (image.Image, error) {
	const maxSide = 8000
	w := int(math.Ceil(float64(c.ink.x1 - c.ink.x0)))
	h := int(math.Ceil(float64(c.ink.y1 - c.ink.y0)))
	if w > maxSide || h > maxSide {
		return nil, fmt.Errorf("drawing too large: %dx%d, at most %d a side", w, h, maxSide)
	}

	ops, clips, bg, hasBG := c.renderPlan()
	tf := ui.Translation(-c.ink.x0, -c.ink.y0)
	img := ui.Render(func(ctx *ui.Context) {
		ctx.Root().Background(ui.Transparent)
		ui.Column(ctx).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			if hasBG {
				p.Fill(ui.Rect{W: float32(w), H: float32(h)}, bg, 0)
			}
			p.Transform(tf, 0, 0, func() { replay(p, ops, clips) })
		})
	}, max(w, 1), max(h, 1), 1)
	return img, nil
}

// renderPlan drops ops before the last Clear and keeps its active clips.
func (c *Context) renderPlan() (ops []drawOp, clips []ui.Path, bg ui.Color, hasBG bool) {
	start := 0
	for i, op := range c.ops {
		if cl, ok := op.(clearOp); ok {
			start, bg, hasBG = i+1, cl.color, true
		}
	}
	for _, op := range c.ops[:start] {
		if cl, ok := op.(clipOp); ok {
			clips = append(clips, cl.path)
		}
	}
	return c.ops[start:], clips, bg, hasBG
}

// replay draws ops in order, applying each clip to everything after it.
func replay(p *ui.Painter, ops []drawOp, clips []ui.Path) {
	var draw func(int)
	draw = func(i int) {
		for ; i < len(ops); i++ {
			if cl, ok := ops[i].(clipOp); ok {
				path := cl.path
				p.ClipPath(&path, func() { draw(i + 1) })
				return
			}
			ops[i].draw(p)
		}
	}
	var wrap func(k int)
	wrap = func(k int) {
		if k >= len(clips) {
			draw(0)
			return
		}
		path := clips[k]
		p.ClipPath(&path, func() { wrap(k + 1) })
	}
	wrap(0)
}
