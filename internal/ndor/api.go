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

func Init() {
	GlobalErr = nil
	GlobalCtx = NewContext()
}
