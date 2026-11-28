package ndor

import (
	"encoding/base64"
	"strings"

	"github.com/zrcoder/ndor/internal/ndor"
	"github.com/zrcoder/ndor/internal/ndor/xgo"
	"github.com/zrcoder/ndor/internal/xgoerr"
)

// Drawing is a recorded ndor drawing, drawn with mygo's painter or exported
// as SVG.
type Drawing = ndor.Drawing

// Compile executes the code and returns its recorded drawing.
func Compile(code string) (*Drawing, *xgoerr.LineError) {
	if strings.TrimSpace(code) == "" {
		return nil, ndor.ErrEmptyInput
	}
	ndor.Init()
	if err := xgo.Run(code); err != nil {
		return nil, err
	}
	if ndor.GlobalErr != nil {
		ndor.GlobalErr.Number = xgoerr.FindLine(code, ndor.GlobalErr.Flag)
		return nil, ndor.GlobalErr
	}
	return ndor.GlobalCtx.Drawing, nil
}

// Run executes the code and returns the SVG as a data URI.
func Run(code string) (string, *xgoerr.LineError) {
	d, lineErr := Compile(code)
	if lineErr != nil {
		return "", lineErr
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(d.SVG()), nil
}
