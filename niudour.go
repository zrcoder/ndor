package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strings"
	"syscall/js"

	"gitee.com/rdor/niudour/internal"
)

var rootDom js.Value

func main() {
	log.SetFlags(log.Lshortfile)
	rootDom = js.Global().Get("window").Get("top")
	rootDom.Get("goButton").Call("addEventListener", "click", js.FuncOf(draw))

	<-make(chan struct{})
}

type requestGopResp struct {
	s   string
	err error
}

func draw(v js.Value, args []js.Value) interface{} {
	rootDom.Get("pictureArea").Set("src", "")
	text := rootDom.Get("getCode").Invoke().String()
	log.Println("receive code:")
	log.Println(text)
	if strings.TrimSpace(text) == "" {
		rootDom.Get("alertError").Invoke(internal.ErrEmptyInput.Error())
		return nil
	}
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

func callGopPlay(input string, c chan<- requestGopResp) {
	log.Println("[test] hello")
	inputData := url.Values{}
	inputData.Add("body", input)
	resp, err := http.DefaultClient.Post(
		"https://play.goplus.org/compile",
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(inputData.Encode()),
	)
	if err != nil {
		log.Println("[test]", err)
		c <- requestGopResp{err: err}
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	gopResp := &struct{ Events []struct{ Message string } }{}
	err = decoder.Decode(gopResp)
	if err != nil {
		log.Println("[test]", err)
		c <- requestGopResp{err: err}
		return
	}
	if len(gopResp.Events) == 0 || gopResp.Events[0].Message == "" {

		log.Println("[test]")

		c <- requestGopResp{err: errors.New("empty response from https://play.goplus.org")}
		return
	}
	c <- requestGopResp{s: gopResp.Events[0].Message}
}
