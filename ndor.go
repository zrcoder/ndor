package ndor

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"

	"github.com/zrcoder/ndor/internal/ndor"
	"github.com/zrcoder/ndor/internal/ndor/xgo"
	"github.com/zrcoder/ndor/internal/xgoerr"
)

// Run executes the code and returns the image as a PNG data URI.
func Run(width, height int, code string) (string, *xgoerr.LineError) {
	src, lineErr := Gen(width, height, code)
	if lineErr != nil {
		return "", lineErr
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(src), nil
}

// Gen executes the code and returns the raw PNG bytes.
func Gen(width, height int, code string) ([]byte, *xgoerr.LineError) {
	if strings.TrimSpace(code) == "" {
		return nil, ndor.ErrEmptyInput
	}
	ndor.Init(width, height)
	if err := xgo.Run(code); err != nil {
		return nil, err
	}
	if ndor.GlobalErr != nil {
		ndor.GlobalErr.Number = xgoerr.FindLine(code, ndor.GlobalErr.Flag)
		return nil, ndor.GlobalErr
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, ndor.GlobalCtx.Image()); err != nil {
		return nil, &xgoerr.LineError{Number: -1, Msg: err.Error()}
	}
	return buf.Bytes(), nil
}
