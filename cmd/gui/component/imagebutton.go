package component

import (
	"github.com/egoist/mygo/ui"
)

// ImageButtonConfig holds configuration for an image button
type ImageButtonConfig struct {
	Width   float32
	Height  float32
	Img     *ui.Bitmap
	OnClick func()
}

// NewImageButton creates an image button with custom config
func NewImageButton(c *ui.Context, cfg ImageButtonConfig) {
	btn := ui.ButtonBase(c).Width(cfg.Width).Height(cfg.Height)
	btn.Background(ui.Color{})
	if cfg.OnClick != nil {
		btn.OnClick(cfg.OnClick)
	}
	btn.Children(func() {
		if cfg.Img != nil {
			w, h := cfg.Width-21, cfg.Height-12
			if btn.Hovered() {
				w *= 0.88
			}
			ui.Image(c, cfg.Img).Width(w).Height(h)
		}
	})
}

// ImageButtonSimple creates an image button with default settings
func ImageButtonSimple(c *ui.Context, img *ui.Bitmap, onClick func()) {
	NewImageButton(c, ImageButtonConfig{
		Width:   ButtonWidth,
		Height:  ButtonHeight,
		Img:     img,
		OnClick: onClick,
	})
}
