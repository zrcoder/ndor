package ndor

import (
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"

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

const textFontFamily = "system-ui, sans-serif"

// pathCmd is one SVG path command, kept as data so a drawing can be replayed
// with mygo's painter and exported as SVG.
type pathCmd struct {
	op  byte // 0 move, 1 line, 2 quad, 3 cubic, 4 close
	pts [3][2]float32
}

type path struct {
	cmds []pathCmd
}

func (p *path) MoveTo(x, y float32) {
	p.cmds = append(p.cmds, pathCmd{op: 0, pts: [3][2]float32{{x, y}}})
}

func (p *path) LineTo(x, y float32) {
	p.cmds = append(p.cmds, pathCmd{op: 1, pts: [3][2]float32{{x, y}}})
}

func (p *path) QuadTo(cx, cy, x, y float32) {
	p.cmds = append(p.cmds, pathCmd{op: 2, pts: [3][2]float32{{cx, cy}, {x, y}}})
}

func (p *path) CubeTo(c1x, c1y, c2x, c2y, x, y float32) {
	p.cmds = append(p.cmds, pathCmd{op: 3, pts: [3][2]float32{{c1x, c1y}, {c2x, c2y}, {x, y}}})
}

func (p *path) Close() {
	p.cmds = append(p.cmds, pathCmd{op: 4})
}

// renderXform scales paths before painting so large drawings do not tile and show seams.
type renderXform struct {
	dx, dy, scale float32
}

func (tf renderXform) x(x float32) float32 { return tf.dx + x*tf.scale }
func (tf renderXform) y(y float32) float32 { return tf.dy + y*tf.scale }

func (p path) ui(tf renderXform) ui.Path {
	var out ui.Path
	for _, c := range p.cmds {
		q := c.pts
		switch c.op {
		case 0:
			out.MoveTo(tf.x(q[0][0]), tf.y(q[0][1]))
		case 1:
			out.LineTo(tf.x(q[0][0]), tf.y(q[0][1]))
		case 2:
			out.QuadTo(tf.x(q[0][0]), tf.y(q[0][1]), tf.x(q[1][0]), tf.y(q[1][1]))
		case 3:
			out.CubeTo(tf.x(q[0][0]), tf.y(q[0][1]), tf.x(q[1][0]), tf.y(q[1][1]), tf.x(q[2][0]), tf.y(q[2][1]))
		case 4:
			out.Close()
		}
	}
	return out
}

func num(v float32) string { return strconv.FormatFloat(float64(v), 'g', -1, 32) }

func (p path) d() string {
	var b strings.Builder
	for i, c := range p.cmds {
		q := c.pts
		if i > 0 {
			b.WriteByte(' ')
		}
		switch c.op {
		case 0:
			fmt.Fprintf(&b, "M%s %s", num(q[0][0]), num(q[0][1]))
		case 1:
			fmt.Fprintf(&b, "L%s %s", num(q[0][0]), num(q[0][1]))
		case 2:
			fmt.Fprintf(&b, "Q%s %s %s %s", num(q[0][0]), num(q[0][1]), num(q[1][0]), num(q[1][1]))
		case 3:
			fmt.Fprintf(&b, "C%s %s %s %s %s %s", num(q[0][0]), num(q[0][1]), num(q[1][0]), num(q[1][1]), num(q[2][0]), num(q[2][1]))
		case 4:
			b.WriteString("Z")
		}
	}
	return b.String()
}

// drawOp is one recorded command, replayed by Render and written by SVG.
type drawOp interface {
	draw(p *ui.Painter, tf renderXform)
	svg(b *strings.Builder)
}

type fillOp struct {
	path  path
	color ui.Color
}

func (o fillOp) draw(p *ui.Painter, tf renderXform) {
	path := o.path.ui(tf)
	p.FillPath(&path, o.color)
}

func (o fillOp) svg(b *strings.Builder) {
	fmt.Fprintf(b, `<path d="%s"`, o.path.d())
	svgPaint(b, "fill", o.color)
	b.WriteString(`/>`)
}

type strokeOp struct {
	path   path
	color  ui.Color
	width  float32
	dashes []float32
}

func (o strokeOp) draw(p *ui.Painter, tf renderXform) {
	path := o.path.ui(tf)
	dashes := make([]float32, len(o.dashes))
	for i, d := range o.dashes {
		dashes[i] = d * tf.scale
	}
	path.Dash(0, dashes...)
	p.StrokePath(&path, o.width*tf.scale, o.color)
}

func (o strokeOp) svg(b *strings.Builder) {
	fmt.Fprintf(b, `<path d="%s" fill="none"`, o.path.d())
	svgPaint(b, "stroke", o.color)
	fmt.Fprintf(b, ` stroke-width="%s" stroke-linecap="round" stroke-linejoin="round"`, num(o.width))
	if len(o.dashes) > 0 {
		b.WriteString(` stroke-dasharray="`)
		for i, d := range o.dashes {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(num(d))
			// SVG repeats an odd-length dash list, as mygo does.
			if i == len(o.dashes)-1 && len(o.dashes)%2 == 1 {
				b.WriteByte(' ')
				b.WriteString(num(d))
			}
		}
		b.WriteByte('"')
	}
	b.WriteString(`/>`)
}

type clipOp struct {
	path path
	id   int
}

func (clipOp) draw(p *ui.Painter, tf renderXform) {}

func (clipOp) svg(b *strings.Builder) {}

type textOp struct {
	text     string
	x, y     float64
	color    ui.Color
	size     float32
	angle    float64
	layout   *ui.TextLayout
	baseline float32
}

func (o textOp) draw(p *ui.Painter, tf renderXform) {
	draw := func() {
		p.TextLayout(o.layout, float32(o.x), float32(o.y)-o.baseline, o.color)
	}
	transform := func() {
		if o.angle != 0 {
			p.Transform(ui.Rotation(float32(o.angle*180/math.Pi)), float32(o.x), float32(o.y), draw)
			return
		}
		draw()
	}
	p.Transform(ui.Translation(tf.dx, tf.dy), 0, 0, func() {
		p.Transform(ui.Scaling(tf.scale, tf.scale), 0, 0, transform)
	})
}

func (o textOp) svg(b *strings.Builder) {
	layout := o.layout.Lines()
	i := 0
	for line := range strings.SplitSeq(o.text, "\n") {
		fmt.Fprintf(b, `<text x="%s"`, num(float32(o.x)))
		y := o.y
		if i < len(layout) {
			y = o.y - float64(o.baseline) + float64(layout[i].Baseline)
		}
		fmt.Fprintf(b, `" y="%s" font-family="%s" font-size="%s" xml:space="preserve"`, num(float32(y)), textFontFamily, num(o.size))
		svgPaint(b, "fill", o.color)
		if o.angle != 0 {
			fmt.Fprintf(b, ` transform="rotate(%s %s %s)"`, num(float32(o.angle*180/math.Pi)), num(float32(o.x)), num(float32(o.y)))
		}
		fmt.Fprintf(b, `>%s</text>`, html.EscapeString(line))
		i++
	}
}

// bounds is the real ink box of the shaped text, rotation included.
func (o textOp) bounds() box {
	w, h := o.layout.Size()
	top := -float64(o.baseline)
	parts := [4][2]float64{{0, top}, {float64(w), top}, {float64(w), top + float64(h)}, {0, top + float64(h)}}
	cos, sin := math.Cos(o.angle), math.Sin(o.angle)
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, q := range parts {
		x, y := o.x+q[0]*cos-q[1]*sin, o.y+q[0]*sin+q[1]*cos
		x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x), max(y1, y)
	}
	return box{x0: float32(x0), y0: float32(y0), x1: float32(x1), y1: float32(y1), ok: true}
}

type clearOp struct {
	color ui.Color
}

func (clearOp) draw(p *ui.Painter, tf renderXform) {}

func (clearOp) svg(b *strings.Builder) {}

func svgPaint(b *strings.Builder, name string, c ui.Color) {
	fmt.Fprintf(b, ` %s="#%02x%02x%02x"`, name, c.R, c.G, c.B)
	if c.A != 255 {
		fmt.Fprintf(b, ` %s-opacity="%s"`, name, num(float32(c.A)/255))
	}
}

// Drawing is a recorded ndor drawing. It holds commands, not a full-size
// bitmap, and can be scaled into any painter rect or exported as SVG.
type Drawing struct {
	ops []drawOp
	ink box
}

// Bounds returns the content box in drawing coordinates.
func (d *Drawing) Bounds() ui.Rect {
	if !d.ink.ok {
		return ui.Rect{}
	}
	return ui.Rect{X: d.ink.x0, Y: d.ink.y0, W: d.ink.x1 - d.ink.x0, H: d.ink.y1 - d.ink.y0}
}

// Render draws the content scaled to fit r, without enlarging it.
func (d *Drawing) Render(p *ui.Painter, r ui.Rect) {
	b := d.Bounds()
	if b.W <= 0 || b.H <= 0 || r.W <= 0 || r.H <= 0 {
		return
	}
	scale := min(r.W/b.W, r.H/b.H, 1)
	tf := renderXform{
		dx:    r.X + (r.W-b.W*scale)/2 - b.X*scale,
		dy:    r.Y + (r.H-b.H*scale)/2 - b.Y*scale,
		scale: scale,
	}
	start, bg, hasBG := d.renderPlan()
	p.Clip(r, 0, func() {
		if hasBG {
			p.Fill(ui.Rect{X: tf.x(b.X), Y: tf.y(b.Y), W: b.W * scale, H: b.H * scale}, bg, 0)
		}
		replay(p, d.ops, start, tf)
	})
}

// SVG returns a standalone SVG document with vector paths and live text.
func (d *Drawing) SVG() []byte {
	b := d.Bounds()
	if b.W <= 0 || b.H <= 0 {
		b = ui.Rect{W: 1, H: 1}
	}
	var out strings.Builder
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="%s %s %s %s">`, num(b.W), num(b.H), num(b.X), num(b.Y), num(b.W), num(b.H))
	start, bg, hasBG := d.renderPlan()
	var defs strings.Builder
	for _, op := range d.ops {
		if cl, ok := op.(clipOp); ok {
			fmt.Fprintf(&defs, `<clipPath id="ndor-clip-%d"><path d="%s"/></clipPath>`, cl.id, cl.path.d())
		}
	}
	if defs.Len() > 0 {
		out.WriteString(`<defs>`)
		out.WriteString(defs.String())
		out.WriteString(`</defs>`)
	}
	if hasBG {
		fmt.Fprintf(&out, `<rect x="%s" y="%s" width="%s" height="%s"`, num(b.X), num(b.Y), num(b.W), num(b.H))
		svgPaint(&out, "fill", bg)
		out.WriteString(`/>`)
	}
	writeSVG(&out, d.ops, start)
	out.WriteString(`</svg>`)
	return []byte(out.String())
}

func writeSVG(b *strings.Builder, ops []drawOp, start int) {
	var draw func(int)
	draw = func(i int) {
		for ; i < len(ops); i++ {
			if cl, ok := ops[i].(clipOp); ok {
				fmt.Fprintf(b, `<g clip-path="url(#ndor-clip-%d)">`, cl.id)
				draw(i + 1)
				b.WriteString(`</g>`)
				return
			}
			if i >= start {
				ops[i].svg(b)
			}
		}
	}
	draw(0)
}

// renderPlan finds the last Clear, whose color becomes the background.
func (d *Drawing) renderPlan() (start int, bg ui.Color, hasBG bool) {
	for i, op := range d.ops {
		if cl, ok := op.(clearOp); ok {
			start, bg, hasBG = i+1, cl.color, true
		}
	}
	return
}

// replay draws ops in order, applying each clip to everything after it.
func replay(p *ui.Painter, ops []drawOp, start int, tf renderXform) {
	var draw func(int)
	draw = func(i int) {
		for ; i < len(ops); i++ {
			if cl, ok := ops[i].(clipOp); ok {
				path := cl.path.ui(tf)
				p.ClipPath(&path, func() { draw(i + 1) })
				return
			}
			if i >= start {
				ops[i].draw(p, tf)
			}
		}
	}
	draw(0)
}

// Context records the ndor API into a Drawing.
type Context struct {
	*Drawing
	matrix      matrix
	fillColor   ui.Color
	strokeColor ui.Color
	textColor   ui.Color
	dashes      []float64
	lineWidth   float64
	stack       []state

	path       path
	hasCurrent bool
	startX     float64
	startY     float64
	curX       float64
	curY       float64

	pathBox box // bounds of the path being built
	clipBox box // bounds of the active clip, if any
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

func intersectBox(a, b box) box {
	if !a.ok || !b.ok {
		return box{}
	}
	x0, y0 := max(a.x0, b.x0), max(a.y0, b.y0)
	x1, y1 := min(a.x1, b.x1), min(a.y1, b.y1)
	if x1 <= x0 || y1 <= y0 {
		return box{}
	}
	return box{x0: x0, y0: y0, x1: x1, y1: y1, ok: true}
}

// NewContext returns a context whose origin is the picture center, +Y up.
func NewContext() *Context {
	return &Context{
		Drawing:     &Drawing{},
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
	c.path = path{}
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
	c.DrawEllipse(x, y, r, r)
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
	size := float32(13 * scale)
	if size < 1 {
		size = 13
	}
	layout := ui.ShapeText(s, ui.Font{Family: textFontFamily, Size: size}, 0)
	lines := layout.Lines()
	baseline := float32(0)
	if len(lines) > 0 {
		baseline = lines[0].Baseline
	}
	op := textOp{text: s, x: tx, y: ty, color: c.textColor, size: size, angle: angle, layout: layout, baseline: baseline}
	c.ops = append(c.ops, op)
	// Text is painted under the active clip; its ink must be clipped too,
	// or an offscreen label enlarges the content box for no drawn pixels.
	r := op.bounds()
	if c.clipBox.ok {
		r = intersectBox(r, c.clipBox)
	}
	if r.ok {
		c.ink.addRect(ui.Rect{X: r.x0, Y: r.y0, W: r.x1 - r.x0, H: r.y1 - r.y0})
	}
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
	c.ops = append(c.ops, clipOp{path: c.path, id: len(c.ops)})
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
