// Package assets loads the images, icons and wiki pages used by the GUI.
package assets

import (
	"io/fs"

	"github.com/egoist/mygo/ui"
)

// Bitmap decodes an embedded PNG from the resources directory.
func Bitmap(fsys fs.FS, name string) *ui.Bitmap {
	data, err := fs.ReadFile(fsys, "resources/"+name)
	if err != nil {
		panic(err)
	}
	b, err := ui.DecodeBitmap(data)
	if err != nil {
		panic(err)
	}
	return b
}

// SVG parses an embedded icon from resources/icons.
func SVG(fsys fs.FS, name string) *ui.SVG {
	data, err := fs.ReadFile(fsys, "resources/icons/"+name)
	if err != nil {
		panic(err)
	}
	return ui.MustParseSVG(data)
}

// Sprite rasterizes an embedded SVG from the resources directory.
func Sprite(fsys fs.FS, name string, size int) *ui.Bitmap {
	data, err := fs.ReadFile(fsys, "resources/"+name)
	if err != nil {
		panic(err)
	}
	return ui.NewBitmap(ui.MustParseSVG(data).Image(size, size, nil))
}

// DocLoader returns a wiki loader for a module: "ndor" or "turtle".
func DocLoader(fsys fs.FS, module string) func(isEnglish bool, tab int) string {
	names := [...]string{"Home.md", "Api.md", "Tutorials.md"}
	dir := "resources/wiki/" + module + "/"

	return func(isEnglish bool, tab int) string {
		if tab < 0 || tab >= len(names) {
			return ""
		}
		lang := ""
		if isEnglish {
			lang = "en/"
		}
		data, _ := fs.ReadFile(fsys, dir+lang+names[tab])
		return string(data)
	}
}
