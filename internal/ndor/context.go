package ndor

import (
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

func identity() matrix {
	return matrix{XX: 1, YY: 1}
}

func translateMatrix(x, y float64) matrix {
	return matrix{XX: 1, YY: 1, X0: x, Y0: y}
}

func scaleMatrix(x, y float64) matrix {
	return matrix{XX: x, YY: y}
}

func rotateMatrix(angle float64) matrix {
	c, s := math.Cos(angle), math.Sin(angle)
	return matrix{XX: c, YX: s, XY: -s, YY: c}
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

func (o textOp) draw(p *ui.Painter) {
	size := float32(13 * o.scale)
	if size < 1 {
		size = 13
	}
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

type clearOp struct {
	color ui.Color
}

func (o clearOp) draw(p *ui.Painter) {}

// Context records the ndor API and renders it with mygo's painter.
type Context struct {
	width, height int

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

	ops []drawOp
}

// NewContext returns a context that draws on a width×height image.
func NewContext(width, height int) *Context {
	return &Context{
		width:       width,
		height:      height,
		matrix:      identity(),
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

func (c *Context) moveTo(x, y float64) {
	x, y = c.matrix.transformPoint(x, y)
	c.path.MoveTo(float32(x), float32(y))
	c.startX, c.startY = x, y
	c.curX, c.curY = x, y
	c.hasCurrent = true
}

func (c *Context) lineTo(x, y float64) {
	if !c.hasCurrent {
		c.moveTo(x, y)
		return
	}
	x, y = c.matrix.transformPoint(x, y)
	c.path.LineTo(float32(x), float32(y))
	c.curX, c.curY = x, y
}

func (c *Context) quadTo(x1, y1, x2, y2 float64) {
	if !c.hasCurrent {
		c.moveTo(x1, y1)
	}
	x1, y1 = c.matrix.transformPoint(x1, y1)
	x2, y2 = c.matrix.transformPoint(x2, y2)
	c.path.QuadTo(float32(x1), float32(y1), float32(x2), float32(y2))
	c.curX, c.curY = x2, y2
}

func (c *Context) cubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !c.hasCurrent {
		c.moveTo(x1, y1)
	}
	x1, y1 = c.matrix.transformPoint(x1, y1)
	x2, y2 = c.matrix.transformPoint(x2, y2)
	x3, y3 = c.matrix.transformPoint(x3, y3)
	c.path.CubeTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(x3), float32(y3))
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
	c.hasCurrent = false
}

func (c *Context) MoveTo(x, y float64) { c.moveTo(x, y) }
func (c *Context) LineTo(x, y float64) { c.lineTo(x, y) }

func (c *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) { c.cubicTo(x1, y1, x2, y2, x3, y3) }
func (c *Context) QuadraticTo(x1, y1, x2, y2 float64)     { c.quadTo(x1, y1, x2, y2) }

func (c *Context) DrawRectangle(x, y, w, h float64) {
	c.newSubPath()
	c.moveTo(x, y)
	c.lineTo(x+w, y)
	c.lineTo(x+w, y+h)
	c.lineTo(x, y+h)
	c.closePath()
}

func (c *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	c.newSubPath()
	c.moveTo(x1, y0)
	c.lineTo(x2, y0)
	c.drawArc(x2, y1, r, Radians(270), Radians(360))
	c.lineTo(x3, y2)
	c.drawArc(x2, y2, r, Radians(0), Radians(90))
	c.lineTo(x1, y3)
	c.drawArc(x1, y2, r, Radians(90), Radians(180))
	c.lineTo(x0, y1)
	c.drawArc(x1, y1, r, Radians(180), Radians(270))
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
				c.lineTo(x0, y0)
			} else {
				c.moveTo(x0, y0)
			}
		}
		c.quadTo(cx, cy, x2, y2)
	}
}

func (c *Context) drawArc(x, y, r, angle1, angle2 float64) {
	c.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

func (c *Context) DrawEllipse(x, y, rx, ry float64) {
	c.newSubPath()
	c.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	c.closePath()
}

func (c *Context) DrawArc(x, y, r, angle1, angle2 float64) { c.drawArc(x, y, r, angle1, angle2) }

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
		c.lineTo(x+r*math.Cos(a), y+r*math.Sin(a))
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
	c.ops = append(c.ops, textOp{text: s, x: tx, y: ty, color: c.textColor, scale: scale, angle: angle})
}

func (c *Context) Fill() {
	c.FillPreserve()
	c.clearPath()
}

func (c *Context) FillPreserve() {
	c.ops = append(c.ops, fillOp{path: c.path, color: c.fillColor})
}

func (c *Context) Stroke() {
	c.StrokePreserve()
	c.clearPath()
}

func (c *Context) StrokePreserve() {
	c.ops = append(c.ops, strokeOp{
		path:   c.path,
		color:  c.strokeColor,
		width:  float32(c.lineWidth),
		dashes: toFloat32s(c.dashes),
	})
}

func (c *Context) Clip() {
	c.ops = append(c.ops, clipOp{path: c.path})
	c.clearPath()
}

func (c *Context) Clear() {
	c.ops = append(c.ops, clearOp{color: c.textColor})
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
	c.matrix = translateMatrix(x, y).multiply(c.matrix)
}

func (c *Context) Scale(x, y float64) {
	c.matrix = scaleMatrix(x, y).multiply(c.matrix)
}

func (c *Context) RotateAbout(angle, x, y float64) {
	c.Translate(x, y)
	c.matrix = rotateMatrix(angle).multiply(c.matrix)
	c.Translate(-x, -y)
}

// Image renders the recorded drawing with mygo's painter and returns it.
func (c *Context) Image() image.Image {
	if c.width <= 0 || c.height <= 0 {
		return image.NewRGBA(image.Rect(0, 0, max(c.width, 0), max(c.height, 0)))
	}
	ops, clips, bg, hasBG := c.renderPlan()
	width, height := c.width, c.height
	return ui.Render(func(ctx *ui.Context) {
		ctx.Root().Background(ui.Transparent)
		ui.Column(ctx).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			if hasBG && bg.A != 0 {
				p.Fill(ui.Rect{W: float32(width), H: float32(height)}, bg, 0)
			}
			replay(p, ops, clips)
		})
	}, width, height, 1)
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

func toFloat32s(vs []float64) []float32 {
	if len(vs) == 0 {
		return nil
	}
	out := make([]float32, len(vs))
	for i, v := range vs {
		out[i] = float32(v)
	}
	return out
}
