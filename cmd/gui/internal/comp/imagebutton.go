package comp

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

const (
	imageButtonWidth  = float32(64)
	imageButtonHeight = float32(72)
)

// ImageButton renders an image button sized for the floating corners. The
// image eases larger while the pointer is over the button.
func ImageButton(c *ui.Context, img *ui.Bitmap, onClick func()) {
	btn := ui.ButtonBase(c).Width(imageButtonWidth).Height(imageButtonHeight).Background(ui.Color{})
	if onClick != nil && btn.Clicked() {
		sound.Play(sound.Tap)
		onClick()
	}
	btn.Children(func() {
		if img == nil {
			return
		}
		image := ui.Image(c, img).Width(imageButtonWidth-21).Height(imageButtonHeight-12).TransformOrigin(0.5, 1)
		p := ui.Pose{}
		if btn.Hovered() {
			p.Scale = 1.2
		}
		image.AnimatePose(p, ui.Snappy)
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
