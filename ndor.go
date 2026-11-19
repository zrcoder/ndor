package ndor

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"github.com/zrcoder/ndor/internal/ndor"
	"github.com/zrcoder/ndor/internal/ndor/xgo"
)

func Run(width, height int, code string) (string, *ndor.LineError) {
	if err := gen(width, height, code); err != nil {
		return "", err
	}
	src, err := getImageSrc(code)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(src), nil
}

func Gen(width, height int, code string) ([]byte, *ndor.LineError) {
	if err := gen(width, height, code); err != nil {
		return nil, err
	}
	return getImageSrc(code)
}

func gen(width, height int, code string) *ndor.LineError {
	if strings.TrimSpace(code) == "" {
		return ndor.ErrEmptyInput
	}
	ndor.Init(width, height)
	return xgo.Run(code)
}

func getImageSrc(oriCode string) ([]byte, *ndor.LineError) {
	if ndor.GlobalErr != nil {
		ndor.GlobalErr.Number = parseErrorline(oriCode, ndor.GlobalErr.Flag)
		return nil, ndor.GlobalErr
	}
	return encode(ndor.GlobalCtx.Image())
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
		line = lowercaseFirstLetter(line)
		if line == flag {
			return j + 1
		}
		j++
	}
	return -1
}

func lowercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func encode(img image.Image) ([]byte, *ndor.LineError) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, &ndor.LineError{Number: -1, Msg: err.Error()}
	}
	return buf.Bytes(), nil
}
