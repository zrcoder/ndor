package ndor

import (
	"github.com/zrcoder/ndor/internal/ndor"
)

const (
	Preserve = "preserve"
)

func Color(args ...any) {
	rgba, err := ndor.ParseColor(args...)
	if err != nil {
		ndor.KeepLastErr(err)
		return
	}
	ndor.GlobalCtx.SetRGBA255(int(rgba.R), int(rgba.G), int(rgba.B), int(rgba.A))
}

func Rectangle(x, y, w, h float64, r ...float64) {
	if len(r) == 0 {
		ndor.GlobalCtx.DrawRectangle(x, y, w, h)
	} else {
		ndor.GlobalCtx.DrawRoundedRectangle(x, y, w, h, r[0])
	}
}

func Circle(x, y, r float64) {
	ndor.GlobalCtx.DrawCircle(x, y, r)
}

func Clip() {
	ndor.GlobalCtx.Clip()
}

func Ellipse(x, y, rx, ry float64) {
	ndor.GlobalCtx.DrawEllipse(x, y, rx, ry)
}

func From(x, y float64) {
	ndor.GlobalCtx.MoveTo(x, y)
}

func To(x, y float64) {
	ndor.GlobalCtx.LineTo(x, y)
}

func Dash(args ...float64) {
	ndor.GlobalCtx.SetDash(args...)
}

func Linew(width float64) {
	ndor.GlobalCtx.SetLineWidth(width)
}

func Text(x, y float64, s string) {
	ndor.GlobalCtx.DrawString(s, x, y)
}

func Arc(x, y, r, angle1, angle2 float64) {
	angle1 = ndor.Radians(angle1)
	angle2 = ndor.Radians(angle2)
	ndor.GlobalCtx.DrawArc(x, y, r, angle1, angle2)
}

func Earc(x, y, rx, ry, angle1, angle2 float64) {
	angle1 = ndor.Radians(angle1)
	angle2 = ndor.Radians(angle2)
	ndor.GlobalCtx.DrawEllipticalArc(x, y, rx, ry, angle1, angle2)
}

func Polygon(n int, x, y, r float64, degree ...float64) {
	deg := 0.0
	if len(degree) > 0 {
		deg = degree[0]
	}
	ndor.GlobalCtx.DrawRegularPolygon(n, x, y, r, ndor.Radians(deg))
}

func Translate(x, y float64) {
	ndor.GlobalCtx.Translate(x, y)
}

func Scale(x, y float64) {
	ndor.GlobalCtx.Scale(x, y)
}

func Rotate(x, y, angle float64) {
	ndor.GlobalCtx.RotateAbout(ndor.Radians(angle), x, y)
}

func Bezier(x1, y1, x2, y2 float64, p ...float64) {
	if len(p) == 2 {
		ndor.GlobalCtx.CubicTo(x1, y1, x2, y2, p[0], p[1])
	} else {
		ndor.GlobalCtx.QuadraticTo(x1, y1, x2, y2)
	}
}

func Fill(preserve ...string) {
	if len(preserve) == 1 && preserve[0] == Preserve {
		ndor.GlobalCtx.FillPreserve()
	} else {
		ndor.GlobalCtx.Fill()
	}
}

func Stroke(preserve ...string) {
	if len(preserve) == 1 && preserve[0] == Preserve {
		ndor.GlobalCtx.StrokePreserve()
	} else {
		ndor.GlobalCtx.Stroke()
	}
}

func Clear() {
	ndor.GlobalCtx.Clear()
}

func Push() {
	ndor.GlobalCtx.Push()
}

func Pop() {
	ndor.GlobalCtx.Pop()
}
