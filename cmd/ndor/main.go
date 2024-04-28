package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zrcoder/ndor/pkg"
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
	src, lerr := pkg.Gen(0, 0, string(data))
	if lerr != nil {
		msg := fmt.Sprintf("line %d: %s", lerr.Number, lerr.Msg)
		log.Fatal(msg)
	}

	outfile := getOutFile(inFile)
	err = os.WriteFile(outfile, src, 0o640)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(outfile, "generated!")
}

func getOutFile(inFile string) string {
	out := strings.TrimSuffix(inFile, filepath.Ext(inFile))
	return out + ".png"
}
