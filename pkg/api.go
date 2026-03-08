package pkg

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"github.com/zrcoder/ndor/internal"
)

func Run(width, height int, code string) (string, *internal.LineError) {
	if err := gen(width, height, code); err != nil {
		return "", err
	}
	src, err := getImageSrc(code)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(src), nil
}

func Gen(width, height int, code string) ([]byte, *internal.LineError) {
	if err := gen(width, height, code); err != nil {
		return nil, err
	}
	return getImageSrc(code)
}

func gen(width, height int, code string) *internal.LineError {
	if strings.TrimSpace(code) == "" {
		return internal.ErrEmptyInput
	}
	internal.Init(width, height)
	return xgoRun(code)
}

func getImageSrc(oriCode string) ([]byte, *internal.LineError) {
	if internal.GlobalErr != nil {
		internal.GlobalErr.Number = parseErrorline(oriCode, internal.GlobalErr.Flag)
		return nil, internal.GlobalErr
	}
	return encode(internal.GlobalCtx.Image())
}

func parseErrorline(oriCode, flag string) int {
	j := 0
	for line := range strings.SplitSeq(oriCode, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			j++
			continue
		}
		i := strings.Index(line, " ")
		if i == -1 {
			i = strings.Index(line, "(")
		}
		if i != -1 {
			line = line[:i]
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

func encode(img image.Image) ([]byte, *internal.LineError) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, &internal.LineError{Number: -1, Msg: err.Error()}
	}
	return buf.Bytes(), nil
}
