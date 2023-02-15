package internal

import (
	"github.com/fogleman/gg"
)

type LineError struct {
	Number int
	Msg    string
}

const (
	invalidPars = "invalid parameters"
)

var (
	GlobalCtx *gg.Context
	GlobalErr *LineError

	ErrEmptyInput = &LineError{
		Number: -1,
		Msg:    "empty input",
	}
)

func Init(width, height int) {
	GlobalErr = nil
	GlobalCtx = gg.NewContext(width, height)
}
