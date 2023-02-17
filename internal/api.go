package internal

import (
	"github.com/fogleman/gg"
)

type LineError struct {
	Number int
	Flag   string
	Msg    string
}

const (
	invalidPars = "invalid parameters"

	colorFlag = "color"
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
