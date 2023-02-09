package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"syscall/js"

	"gitee.com/rdor/niudour/internal"
)

var rootDom js.Value

func main() {
	log.SetFlags(log.Lshortfile)
	rootDom = js.Global().Get("window").Get("top")
	rootDom.Get("goButton").Call("addEventListener", "click", js.FuncOf(draw))
	select {}
}

func draw(v js.Value, args []js.Value) interface{} {
	rootDom.Get("pictureArea").Set("src", "")
	text := rootDom.Get("getCode").Invoke().String()
	log.Println("receive code:")
	log.Println(text)
	painter := internal.NewPainter(text)
	pictureBox := rootDom.Get("pictureBox")
	w := pictureBox.Get("offsetWidth").Int()
	h := pictureBox.Get("offsetHeight").Int()
	img, err := painter.Draw(w, h)
	if err != nil {
		log.Println(err)
		rootDom.Get("alertError").Invoke(err.Error())
		return nil
	}
	src, err := getImgString(img)
	if err != nil {
		rootDom.Get("alertError").Invoke(err.Error())
		return nil
	}
	rootDom.Get("toastSuccess").Invoke(src)
	return nil
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
