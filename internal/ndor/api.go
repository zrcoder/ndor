package ndor

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
	GlobalCtx *Context
	GlobalErr *LineError

	ErrEmptyInput = &LineError{
		Number: -1,
		Msg:    "empty input",
	}
)

func Init(width, height int) {
	GlobalErr = nil
	GlobalCtx = NewContext(width, height)
}
