package ndor

import (
	icolor "image/color"

	"github.com/zrcoder/ndor/internal/parse"
)

func ParseColor(args ...any) (icolor.RGBA, *LineError) {
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

func KeepLastErr(err *LineError) {
	if GlobalErr != nil {
		return
	}
	GlobalErr = err
}

func newInvalidParsErr(flag, msg string) *LineError {
	return &LineError{
		Flag: flag,
		Msg:  msg,
	}
}
