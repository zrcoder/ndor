package comp

import "github.com/egoist/mygo/ui"

const (
	imageButtonWidth  = float32(64)
	imageButtonHeight = float32(72)
)

// ImageButton renders an image button sized for the floating corners.
func ImageButton(c *ui.Context, img *ui.Bitmap, onClick func()) {
	btn := ui.ButtonBase(c).Width(imageButtonWidth).Height(imageButtonHeight).Background(ui.Color{})
	if onClick != nil {
		btn.OnClick(onClick)
	}
	btn.Children(func() {
		if img == nil {
			return
		}
		w, h := imageButtonWidth-21, imageButtonHeight-12
		if btn.Hovered() {
			w *= 0.88
		}
		ui.Image(c, img).Width(w).Height(h)
	})
}

// FloatButtons renders the help and run buttons in the bottom corners.
func FloatButtons(c *ui.Context, helpImg, runImg *ui.Bitmap, onHelp, onRun func()) {
	ui.Column(c.Key("help-float")).Width(imageButtonWidth).Height(imageButtonHeight).Bottom(0).Left(4).Absolute().Children(func() {
		ImageButton(c, helpImg, onHelp)
	})
	ui.Column(c.Key("run-float")).Width(imageButtonWidth).Height(imageButtonHeight).Bottom(0).Right(0).Absolute().Children(func() {
		ImageButton(c, runImg, onRun)
	})
}
