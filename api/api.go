package api

import (
	"github.com/zrcoder/ndor/internal"

	"github.com/fogleman/gg"
)

const (
	Preserve = "preserve"
)

func Context(w, h int) {
	internal.GlobalCtx = gg.NewContext(w, h)
}

func Color(args ...any) {
	rgba, err := internal.ParseColor(args...)
	if err != nil {
		internal.KeepLastErr(err)
		return
	}
	internal.GlobalCtx.SetRGBA255(int(rgba.R), int(rgba.G), int(rgba.B), int(rgba.A))
}

func Rectangle(x, y, w, h float64, r ...float64) {
	if len(r) == 0 {
		internal.GlobalCtx.DrawRectangle(x, y, w, h)
	} else {
		internal.GlobalCtx.DrawRoundedRectangle(x, y, w, h, r[0])
	}
}

func Circle(x, y, r float64) {
	internal.GlobalCtx.DrawCircle(x, y, r)
}

func Clip() {
	internal.GlobalCtx.Clip()
}

func Ellipse(x, y, rx, ry float64) {
	internal.GlobalCtx.DrawEllipse(x, y, rx, ry)
}

func From(x, y float64) {
	internal.GlobalCtx.MoveTo(x, y)
}

func To(x, y float64) {
	internal.GlobalCtx.LineTo(x, y)
}

func Dash(args ...float64) {
	internal.GlobalCtx.SetDash(args...)
}

func Linew(width float64) {
	internal.GlobalCtx.SetLineWidth(width)
}

func Text(x, y float64, s string) {
	internal.GlobalCtx.DrawString(s, x, y)
}

func Arc(x, y, r, angle1, angle2 float64) {
	angle1 = gg.Radians(angle1)
	angle2 = gg.Radians(angle2)
	internal.GlobalCtx.DrawArc(x, y, r, angle1, angle2)
}

func Earc(x, y, rx, ry, angle1, angle2 float64) {
	angle1 = gg.Radians(angle1)
	angle2 = gg.Radians(angle2)
	internal.GlobalCtx.DrawEllipticalArc(x, y, rx, ry, angle1, angle2)
}

func Polygon(n int, x, y, r float64, degree ...float64) {
	deg := 0.0
	if len(degree) > 0 {
		deg = degree[0]
	}
	internal.GlobalCtx.DrawRegularPolygon(n, x, y, r, gg.Radians(deg))
}

func Translate(x, y float64) {
	internal.GlobalCtx.Translate(x, y)
}

func Scale(x, y float64) {
	internal.GlobalCtx.Scale(x, y)
}

func Rotate(x, y, angle float64) {
	internal.GlobalCtx.RotateAbout(gg.Radians(angle), x, y)
}

func Bezier(x1, y1, x2, y2 float64, p ...float64) {
	if len(p) == 2 {
		internal.GlobalCtx.CubicTo(x1, y1, x2, y2, p[0], p[1])
	} else {
		internal.GlobalCtx.QuadraticTo(x1, y1, x2, y2)
	}
}

func Fill(preserve ...string) {
	if len(preserve) == 1 && preserve[0] == Preserve {
		internal.GlobalCtx.FillPreserve()
	} else {
		internal.GlobalCtx.Fill()
	}
}

func Stroke(preserve ...string) {
	if len(preserve) == 1 && preserve[0] == Preserve {
		internal.GlobalCtx.StrokePreserve()
	} else {
		internal.GlobalCtx.Stroke()
	}
}

func Clear() {
	internal.GlobalCtx.Clear()
}

func Push() {
	internal.GlobalCtx.Push()
}

func Pop() {
	internal.GlobalCtx.Pop()
}
