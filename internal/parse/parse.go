// Package parse converts xgo arguments into colors and floats.
package parse

import (
	"errors"
	icolor "image/color"
	"strconv"

	"golang.org/x/image/colornames"
)

// Color parses a color name, a hex string, or 3-4 numeric components.
func Color(args ...any) (icolor.RGBA, bool) {
	switch len(args) {
	case 1:
		s, ok := args[0].(string)
		if !ok {
			return icolor.RGBA{}, false
		}
		if c, ok := colornames.Map[s]; ok {
			return c, true
		}
		return hexColor(s)
	case 3, 4:
		floats, err := Floats(args)
		if err != nil {
			return icolor.RGBA{}, false
		}
		var a uint8 = 0xff
		if len(args) == 4 {
			a = uint8(floats[3])
		}
		return icolor.RGBA{
			R: uint8(floats[0]),
			G: uint8(floats[1]),
			B: uint8(floats[2]),
			A: a,
		}, true
	}
	return icolor.RGBA{}, false
}

// Floats converts each argument to a float64.
func Floats(pars []any) ([]float64, error) {
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

func hexColor(s string) (icolor.RGBA, bool) {
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
	return 0, errors.New("invalid parameter")
}
