package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/zrcoder/ndor"
	nmd "github.com/zrcoder/ndor/internal/goldmark"
)

func main() {
	log.SetFlags(log.LUTC)
	if len(os.Args) < 2 {
		log.Fatal("must input the source file")
	}
	inFile := os.Args[1]
	data, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}
	outFile := ""
	if strings.HasSuffix(inFile, ".md") {
		outFile = getOutfile(inFile, ".html")
		md := goldmark.New(goldmark.WithExtensions(&nmd.Extender{}))
		buf := new(bytes.Buffer)
		err = md.Convert(data, buf)
		data = buf.Bytes()
	} else {
		outFile = getOutfile(inFile, ".png")
		data, err = genPng(data)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(outFile, data, 0o640)
	if err != nil {
		log.Fatal(err)
	}
}

func genPng(data []byte) ([]byte, error) {
	src, lerr := ndor.Gen(string(data))
	if lerr != nil {
		return nil, fmt.Errorf("line %d: %s", lerr.Number, lerr.Msg)
	}
	return src, nil
}

func getOutfile(inFile, ext string) string {
	return strings.TrimSuffix(inFile, filepath.Ext(inFile)) + ext
}
