package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zrcoder/niudour/pkg"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must input the source file")
	}
	inFile := os.Args[1]
	data, err := os.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}
	src, lerr := pkg.Gen(0, 0, string(data))
	if lerr != nil {
		log.Fatal(lerr)
	}

	outfile := getOutFile(inFile)
	err = os.WriteFile(outfile, src, 0640)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(outfile, "generated!")
}

func getOutFile(inFile string) string {
	out := strings.TrimSuffix(inFile, filepath.Ext(inFile))
	return out + ".png"
}
