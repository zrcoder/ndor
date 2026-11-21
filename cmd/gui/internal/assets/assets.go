// Package assets loads the images, icons and wiki pages used by the GUI.
package assets

import (
	"bytes"
	"image"
	"image/png"
	"io/fs"

	"github.com/egoist/mygo/ui"
	"golang.org/x/image/draw"
)

func read(fsys fs.FS, name string) []byte {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		panic(err)
	}
	return data
}

// Bitmap decodes an embedded PNG from the resources directory.
func Bitmap(fsys fs.FS, name string) *ui.Bitmap {
	b, err := ui.DecodeBitmap(read(fsys, "resources/"+name))
	if err != nil {
		panic(err)
	}
	return b
}

// SVG parses an embedded icon from resources/icons.
func SVG(fsys fs.FS, name string) *ui.SVG {
	return ui.MustParseSVG(read(fsys, "resources/icons/"+name))
}

// SpriteSheet decodes a horizontal sprite sheet into per-frame bitmaps.
func SpriteSheet(fsys fs.FS, name string, frames, size int) []*ui.Bitmap {
	img, err := png.Decode(bytes.NewReader(read(fsys, "resources/"+name)))
	if err != nil {
		panic(err)
	}
	b := img.Bounds()
	fw := b.Dx() / frames
	out := make([]*ui.Bitmap, frames)
	for i := range out {
		frame := image.NewRGBA(image.Rect(0, 0, size, size))
		src := image.Rect(b.Min.X+i*fw, b.Min.Y, b.Min.X+(i+1)*fw, b.Max.Y)
		draw.CatmullRom.Scale(frame, frame.Bounds(), img, src, draw.Over, nil)
		out[i] = ui.NewBitmap(frame)
	}
	return out
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
