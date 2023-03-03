package internal

import (
	"errors"
	icolor "image/color"
	"strconv"

	"golang.org/x/image/colornames"
)

func ParseColor(args ...any) (icolor.RGBA, *LineError) {
	var rgba icolor.RGBA
	switch len(args) {
	case 1:
		s, ok := args[0].(string)
		if !ok {
			return rgba, newInvalidParsErr(colorFlag, invalidPars)
		}
		if c, ok := colornames.Map[s]; ok {
			rgba = c
		} else {
			if c, ok = parseHexColor(s); !ok {
				return rgba, newInvalidParsErr(colorFlag, "invalid color name or hex color format")
			}
			rgba = c
		}
	case 3, 4:
		floats, err := toFloats((args))
		if err != nil {
			return rgba, newInvalidParsErr(colorFlag, invalidPars)
		}
		r, g, b := floats[0], floats[1], floats[2]
		var a uint8 = 0xff
		if len(args) == 4 {
			a = uint8(floats[3])
		}
		rgba = icolor.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: a}
	default:
		return rgba, newInvalidParsErr(colorFlag, invalidPars)
	}
	return rgba, nil
}

func KeepLastErr(err *LineError) {
	if GlobalErr != nil {
		return
	}
	GlobalErr = err
}

func parseHexColor(s string) (icolor.RGBA, bool) {
	if s == "" || s[0] != '#' {
		return icolor.RGBA{}, false
	}
	res := true
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		default:
			res = false
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
		res = false
	}
	return c, res
}

func toFloats(pars []any) ([]float64, error) {
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

func parseFloat(param any) (float64, error) {
	switch v := param.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, errors.New(invalidPars)
}

func newInvalidParsErr(flag, msg string) *LineError {
	return &LineError{
		Flag: flag,
		Msg:  msg,
	}
}
