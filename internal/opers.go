package internal

import (
	icolor "image/color"
	"math"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
)

func context(p *Painter, line lineInfo) error {
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 2 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	w, h := floats[0], floats[1]
	p.context = gg.NewContext(int(w), int(h))
	return nil
}

func color(p *Painter, line lineInfo) error {
	ensureContext(p)
	n := len(line.Params)
	invalidErr := genInvalidParamError(line.Num)
	var rgba icolor.RGBA
	switch n {
	case 1:
		if c, ok := colornames.Map[line.Params[0]]; ok {
			rgba = c
		} else {
			c, err := parseHexColor(line.Params[0], invalidErr)
			if err != nil {
				return err
			}
			rgba = c
		}
	case 3, 4:
		floats, err := toFloats((line.Params))
		if err != nil {
			return invalidErr
		}
		r, g, b := floats[0], floats[1], floats[2]
		var a uint8 = 0xff
		if n == 4 {
			a = uint8(floats[3] * float64(a))
		}
		rgba = icolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a}
	default:
		return invalidErr
	}
	p.context.SetRGBA255(int(rgba.R), int(rgba.G), int(rgba.B), int(rgba.A))
	return nil
}

func parseHexColor(s string, invalidErr error) (icolor.RGBA, error) {
	if s == "" || s[0] != '#' {
		return icolor.RGBA{}, invalidErr
	}
	var err error
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		default:
			err = invalidErr
			return 0
		}
	}
	c := icolor.RGBA{A: 0xff}
	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = invalidErr
	}
	return c, err
}

func rectanger(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 4 && len(line.Params) != 5 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y, w, h := floats[0], floats[1], floats[2], floats[3]
	if len(line.Params) == 4 {
		p.context.DrawRectangle(x, y, w, h)
		return nil
	}
	r := floats[4]
	p.context.DrawRoundedRectangle(x, y, w, h, r)
	return nil
}

func circle(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y, r := floats[0], floats[1], floats[2]
	p.context.DrawCircle(x, y, r)
	return nil
}

func clip(p *Painter, line lineInfo) error {
	ensureContext(p)
	p.context.Clip()
	return nil
}

func ellipse(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 4 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y, rx, ry := floats[0], floats[1], floats[2], floats[3]
	p.context.DrawEllipse(x, y, rx, ry)
	return nil
}

func moveTo(p *Painter, line lineInfo) error {
	ensureContext(p)
	return dealWithXY(line, p.context.MoveTo)
}

func lineTo(p *Painter, line lineInfo) error {
	ensureContext(p)
	return dealWithXY(line, p.context.LineTo)
}

func dealWithXY(line lineInfo, f func(x, y float64)) error {
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 2 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y := floats[0], floats[1]
	f(x, y)
	return nil
}

func dash(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	p.context.SetDash(floats...)
	return nil
}

func setLineWidth(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 1 {
		return invalidErr
	}
	w, err := parseFloat(line.Params[0])
	if err != nil {
		return invalidErr
	}
	p.context.SetLineWidth(w)
	return nil
}

func text(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) < 3 {
		return invalidErr
	}
	x, err := parseFloat(line.Params[0])
	if err != nil {
		return invalidErr
	}
	y, err := parseFloat(line.Params[1])
	if err != nil {
		return invalidErr
	}
	str := strings.Join(line.Params[2:], " ")
	p.context.DrawString(str, x, y)
	return nil
}

func arc(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 5 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y, r := floats[0], floats[1], floats[2]
	angle1, angle2 := radians(floats[3]), radians(floats[4])
	p.context.DrawArc(x, y, r, angle1, angle2)
	return nil
}

func ellipseArc(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 6 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y, rx, ry := floats[0], floats[1], floats[2], floats[3]
	angle1, angle2 := radians(floats[4]), radians(floats[5])
	p.context.DrawEllipticalArc(x, y, rx, ry, angle1, angle2)
	return nil
}

func regularPolygon(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) < 4 || len(line.Params) > 5 {
		return invalidErr
	}
	n, err := strconv.Atoi(line.Params[0])
	if err != nil {
		return invalidErr
	}
	floats, err := toFloats(line.Params[1:])
	if err != nil {
		return invalidErr
	}
	x, y, r := floats[0], floats[1], floats[2]
	degree := 0.0
	if len(line.Params) == 5 {
		degree = radians(floats[3])
	}
	p.context.DrawRegularPolygon(n, x, y, r, degree)
	return nil
}

func translate(p *Painter, line lineInfo) error {
	ensureContext(p)
	return dealWithXY(line, p.context.Translate)
}

func scale(p *Painter, line lineInfo) error {
	ensureContext(p)
	return dealWithXY(line, p.context.Scale)
}

func rotate(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 3 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x, y := floats[0], floats[1]
	angle := radians(floats[2])
	p.context.RotateAbout(angle, x, y)
	return nil
}

func bezier(p *Painter, line lineInfo) error {
	ensureContext(p)
	invalidErr := genInvalidParamError(line.Num)
	if len(line.Params) != 4 && len(line.Params) != 6 {
		return invalidErr
	}
	floats, err := toFloats(line.Params)
	if err != nil {
		return invalidErr
	}
	x1, y1, x2, y2 := floats[0], floats[1], floats[2], floats[3]
	if len(line.Params) == 4 {
		p.context.QuadraticTo(x1, y1, x2, y2)
		return nil
	}
	x3, y3 := floats[4], floats[5]
	p.context.CubicTo(x1, y1, x2, y2, x3, y3)
	return nil
}

func fill(p *Painter, line lineInfo) error {
	ensureContext(p)
	if len(line.Params) == 0 {
		p.context.Fill()
		return nil
	}
	if len(line.Params) != 1 || line.Params[0] != "preserve" {
		return genInvalidParamError(line.Num)
	}
	p.context.FillPreserve()
	return nil
}

func stroke(p *Painter, line lineInfo) error {
	ensureContext(p)
	if len(line.Params) == 0 {
		p.context.Stroke()
		return nil
	}
	if len(line.Params) != 1 || line.Params[0] != "preserve" {
		return genInvalidParamError(line.Num)
	}
	p.context.StrokePreserve()
	return nil
}

func clear(p *Painter, line lineInfo) error {
	ensureContext(p)
	p.context.Clear()
	return nil
}

func push(p *Painter, line lineInfo) error {
	ensureContext(p)
	p.context.Push()
	return nil
}

func pop(p *Painter, line lineInfo) error {
	ensureContext(p)
	p.context.Pop()
	return nil
}

func toFloats(pars []string) ([]float64, error) {
	res := make([]float64, len(pars))
	for i, s := range pars {
		n, err := parseFloat(s)
		if err != nil {
			return nil, err
		}
		res[i] = n
	}
	return res, nil
}

func genInvalidParamError(num int) error {
	return genError(num, "invalid parameters")
}

func genError(num int, msg string) error {
	return &LineError{
		Msg:    msg,
		Number: num,
	}
}

const defaultContextSize = 1024

func ensureContext(p *Painter) {
	if p.context == nil {
		p.context = gg.NewContext(defaultContextSize, defaultContextSize)
	}
}

func parseFloat(param string) (float64, error) {
	return strconv.ParseFloat(param, 64)
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
