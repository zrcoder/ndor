//go:build ignore

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func main() {
	if err := generateHelpImage(); err != nil {
		log.Fatal(err)
	}
	if err := generateGoImage(); err != nil {
		log.Fatal(err)
	}
}

func loadFont(size float64) (*truetype.Font, font.Face) {
	fontPath := "/System/Library/Fonts/Supplemental/Arial Rounded Bold.ttf"

	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("failed to load font: %v", err)
	}

	f, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	// Use no hinting and lower DPI for simpler rendering
	face := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	return f, face
}

func measureString(face font.Face, text string) int {
	advance := fixed.I(0)
	for _, r := range text {
		adv, ok := face.GlyphAdvance(r)
		if ok {
			advance += adv
		}
	}
	return advance.Ceil()
}

func drawText(img *image.RGBA, face font.Face, text string, x, y int, textColor color.RGBA) {
	// Draw text with antialiasing
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(text)
}

func generateHelpImage() error {
	f, err := os.Open("resources/teacher.png")
	if err != nil {
		return err
	}
	defer f.Close()

	orig, err := png.Decode(f)
	if err != nil {
		return err
	}

	// Create a new image with the same config as orig, then draw text on it
	bounds := orig.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, orig, bounds.Min, draw.Over)

	w := bounds.Dx()
	text := "HELP"

	_, face := loadFont(81)
	textW := measureString(face, text)
	textX := (w - textW) / 2
	textY := bounds.Dy() - 20

	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	drawText(img, face, text, textX, textY, black)

	out, err := os.Create("resources/teacher_labeled.png")
	if err != nil {
		return err
	}
	defer out.Close()

	// Use optimal compression
	return png.Encode(out, img)
}

func generateGoImage() error {
	f, err := os.Open("resources/gopher.png")
	if err != nil {
		return err
	}
	defer f.Close()

	orig, err := png.Decode(f)
	if err != nil {
		return err
	}

	bounds := orig.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, orig, bounds.Min, draw.Over)

	w := bounds.Dx()
	text := "GO"

	_, face := loadFont(81)
	textW := measureString(face, text)
	textX := (w - textW) / 2
	textY := bounds.Dy() - 20

	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	drawText(img, face, text, textX, textY, black)

	out, err := os.Create("resources/gopher_labeled.png")
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, img)
}
