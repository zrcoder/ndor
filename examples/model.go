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
	//go:embed *.xgo
	fs      embed.FS
	Default []Example
)

func init() {
	files := []string{
		"hello-world.xgo",
		"hello-word-go.xgo",
		"beziel.xgo",
		"clip.xgo",
		"crisp.xgo",
		"fan.xgo",
		"flower.xgo",
		"geometry.xgo",
		"line-width.xgo",
		"lines.xgo",
		"open-fill.xgo",
		"rainbow.xgo",
		"spiral.xgo",
		"star.xgo",
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
