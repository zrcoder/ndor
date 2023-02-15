package internal

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"strings"
)

func Draw(width, height int, code string) (imgSrc string, errLine int, errInfo string) {
	log.Printf("origin code to run:\n%s\n", code)
	if strings.TrimSpace(code) == "" {
		return "", -1, ErrEmptyInput.Error()
	}
	var err error
	preLines := strings.Count(preCodes, "\n")
	code, err = gopExecute(preCodes+code, preLines)
	if err != nil {
		if lineErr, ok := err.(*LineError); ok {
			return "", lineErr.Number, lineErr.Msg
		}
		return "", -1, err.Error()
	}
	log.Printf("code after gop exec:\n%s\n", code)

	painter := NewPainter(code, preLines)
	img, err := painter.Draw(width, height)
	if err != nil {
		if lineErr, ok := err.(*LineError); ok {
			return "", lineErr.Number, lineErr.Msg
		}
		return "", -1, err.Error()
	}

	return getImgString(img)
}

func getImgString(img image.Image) (string, int, string) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return "", -1, err.Error()
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), -1, ""
}
