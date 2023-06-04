package examples

import (
	"embed"
	"log"
	"path/filepath"
	"strings"
)

type Example struct {
	Name string
	Code string
}

var (
	//go:embed *.gop
	fs      embed.FS
	Default []Example
)

func init() {
	files := []string{
		"hello-world.gop",
		"hello-word-go.gop",
		"beziel.gop",
		"clip.gop",
		"crisp.gop",
		"fan.gop",
		"flower.gop",
		"geometry.gop",
		"line-width.gop",
		"lines.gop",
		"open-fill.gop",
		"rainbow.gop",
		"spiral.gop",
		"start.gop",
	}
	Default = make([]Example, 0, len(files))

	for _, n := range files {
		data, err := fs.ReadFile(n)
		if err != nil {
			log.Fatal(err)
		}
		name := strings.TrimSuffix(n, filepath.Ext(n))
		code := string(data)
		i := strings.Index(code, "\n")
		if i > 0 {
			name = code[:i]
			name = strings.TrimLeft(name, "/ ")
		}
		Default = append(Default, Example{name, code})
	}
}
