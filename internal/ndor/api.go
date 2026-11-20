package ndor

import "github.com/zrcoder/ndor/internal/xgoerr"

const (
	invalidPars = "invalid parameters"

	colorFlag = "color"
)

var (
	GlobalCtx *Context
	GlobalErr *xgoerr.LineError

	ErrEmptyInput = &xgoerr.LineError{
		Number: -1,
		Msg:    "empty input",
	}
)

func Init(width, height int) {
	GlobalErr = nil
	GlobalCtx = NewContext(width, height)
}
