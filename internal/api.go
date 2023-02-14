package internal

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"strings"
)

func Draw(width, height int, code string) (string, error) {
	log.Println("[test]")
	if strings.TrimSpace(code) == "" {
		return "", errEmptyInput
	}

	var err error

	// code, err = gopExecute(regularInput(code))
	// if err != nil {
	// 	return "", err
	// }

	painter := NewPainter(code)
	img, err := painter.Draw(width, height)
	if err != nil {
		return "", err
	}
	return getImgString(img)
}

func getImgString(img image.Image) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
