package ndor

import (
	icolor "image/color"

	"github.com/zrcoder/ndor/internal/parse"
	"github.com/zrcoder/ndor/internal/xgoerr"
)

func ParseColor(args ...any) (icolor.RGBA, *xgoerr.LineError) {
	rgba, ok := parse.Color(args...)
	if ok {
		return rgba, nil
	}
	if len(args) == 1 {
		if _, isStr := args[0].(string); isStr {
			return rgba, newInvalidParsErr(colorFlag, "invalid color name or hex color format")
		}
	}
	return rgba, newInvalidParsErr(colorFlag, invalidPars)
}

func newInvalidParsErr(flag, msg string) *xgoerr.LineError {
	return &xgoerr.LineError{
		Flag: flag,
		Msg:  msg,
	}
}
