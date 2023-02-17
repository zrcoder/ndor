package pkg

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"strings"

	"github.com/zrcoder/niudour/internal"
)

func Run(width, height int, code string) (string, *internal.LineError) {
	if strings.TrimSpace(code) == "" {
		return "", internal.ErrEmptyInput
	}

	internal.Init(width, height)
	err := gopRun(code)
	if err != nil {
		return "", err
	}
	return getImageSrc(code)
}

func getImageSrc(oriCode string) (string, *internal.LineError) {
	if internal.GlobalErr != nil {
		internal.GlobalErr.Number = parseErrorline(oriCode, internal.GlobalErr.Flag)
		return "", internal.GlobalErr
	}
	return encode(internal.GlobalCtx.Image())
}

func parseErrorline(oriCode, flag string) int {
	arr := strings.Split(oriCode, "\n")
	for j, line := range arr {
		line = strings.TrimSpace(line)
		if line == "" {
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
	}
	return -1
}

func lowercaseFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func encode(img image.Image) (string, *internal.LineError) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return "", &internal.LineError{Number: -1, Msg: err.Error()}
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
