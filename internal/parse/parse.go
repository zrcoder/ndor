// Package parse converts xgo arguments into colors and floats.
package parse

import (
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
		values := [4]float64{0, 0, 0, 0xff}
		for i, arg := range args {
			f, ok := toFloat(arg)
			if !ok {
				return icolor.RGBA{}, false
			}
			values[i] = f
		}
		return icolor.RGBA{R: uint8(values[0]), G: uint8(values[1]), B: uint8(values[2]), A: uint8(values[3])}, true
	}
	return icolor.RGBA{}, false
}

func hexColor(s string) (icolor.RGBA, bool) {
	if s == "" || s[0] != '#' {
		return icolor.RGBA{}, false
	}
	switch len(s) {
	case 7:
		v, err := strconv.ParseUint(s[1:], 16, 32)
		if err != nil {
			return icolor.RGBA{}, false
		}
		return icolor.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}, true
	case 4:
		v, err := strconv.ParseUint(s[1:], 16, 16)
		if err != nil {
			return icolor.RGBA{}, false
		}
		return icolor.RGBA{R: uint8(v>>8&0xf) * 17, G: uint8(v>>4&0xf) * 17, B: uint8(v&0xf) * 17, A: 0xff}, true
	}
	return icolor.RGBA{}, false
}

func toFloat(param any) (float64, bool) {
	switch v := param.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return f, err == nil
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}
