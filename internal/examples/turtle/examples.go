package examples

import (
	"embed"
	"log"
	"strings"
)

// Examples represents an example with names and codes
type Examples map[string]string

var (
	//go:embed *.xgo
	fs embed.FS

	// Default is the list of all available examples
	Default Examples

	DefaultNames = []string{
		"star",
		"tree",
		"colored squares",
		"flower",
		"snowflake",
		"sierpinski",
		"fibonacci",
	}
)

func init() {
	Default = make(Examples, len(DefaultNames))

	for _, n := range DefaultNames {
		path := strings.ReplaceAll(n, " ", "-") + ".xgo"
		data, err := fs.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		Default[n] = string(data)
	}
}
