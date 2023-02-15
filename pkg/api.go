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
	return getImageSrc()
}

func getImageSrc() (string, *internal.LineError) {
	if internal.GlobalErr != nil {
		return "", internal.GlobalErr
	}
	return encode(internal.GlobalCtx.Image())
}

func encode(img image.Image) (string, *internal.LineError) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return "", &internal.LineError{Number: -1, Msg: err.Error()}
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
