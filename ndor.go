package ndor

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"

	"github.com/zrcoder/ndor/internal/ndor"
	"github.com/zrcoder/ndor/internal/ndor/xgo"
)

// Run executes the code and returns the image as a PNG data URI.
func Run(width, height int, code string) (string, *ndor.LineError) {
	src, lineErr := Gen(width, height, code)
	if lineErr != nil {
		return "", lineErr
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(src), nil
}

// Gen executes the code and returns the raw PNG bytes.
func Gen(width, height int, code string) ([]byte, *ndor.LineError) {
	if strings.TrimSpace(code) == "" {
		return nil, ndor.ErrEmptyInput
	}
	ndor.Init(width, height)
	if err := xgo.Run(code); err != nil {
		return nil, err
	}
	if ndor.GlobalErr != nil {
		ndor.GlobalErr.Number = parseErrorline(code, ndor.GlobalErr.Flag)
		return nil, ndor.GlobalErr
	}
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, ndor.GlobalCtx.Image()); err != nil {
		return nil, &ndor.LineError{Number: -1, Msg: err.Error()}
	}
	return buf.Bytes(), nil
}

func parseErrorline(oriCode, flag string) int {
	j := 0
	for line := range strings.SplitSeq(oriCode, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			j++
			continue
		}
		before, _, found := strings.Cut(line, " ")
		if !found {
			before, _, found = strings.Cut(line, "(")
		}
		if found {
			line = before
		}
		if strings.ToLower(line[:1])+line[1:] == flag {
			return j + 1
		}
		j++
	}
	return -1
}
