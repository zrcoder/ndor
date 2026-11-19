package comp

import (
	"github.com/egoist/mygo/ui"
)

type DialogOption struct {
	Open    *bool
	Img     *ui.Bitmap
	Buttons []string
	OnClick func(btnIndex int)
}

// ShowDialogWithImage shows an image dialog; the last button is the primary action.
func ShowDialogWithImage(c *ui.Context, opt DialogOption) {
	if !*opt.Open {
		return
	}

	if ui.Modal(c, opt.Open, func() {
		ui.Column(c).Width(250).Radius(12).Background(c.Theme().Surface).Children(func() {
			if opt.Img != nil {
				ui.Row(c).Fill().Children(func() {
					ui.Spacer(c)
					ui.Image(c, opt.Img).Width(80).Height(100)
					ui.Spacer(c)
				})
			}

			if len(opt.Buttons) > 0 {
				ui.Row(c).Fill().Children(func() {
					ui.Spacer(c)
					for i, btnText := range opt.Buttons {
						idx := i
						bgColor := c.Theme().SurfaceHover
						textColor := c.Theme().Text
						if i == len(opt.Buttons)-1 {
							bgColor = c.Theme().Accent
							textColor = c.Theme().AccentText
						}
						ui.ButtonBase(c).Width(100).Height(36).Radius(6).Background(bgColor).OnClick(func() {
							*opt.Open = false
							if opt.OnClick != nil {
								opt.OnClick(idx)
							}
						}).Children(func() {
							ui.Text(c, btnText).FontSize(14).TextColor(textColor)
						})
						if i < len(opt.Buttons)-1 {
							ui.Spacer(c).Width(8)
						}
					}
					ui.Spacer(c)
				})
			}
		})
	}).Dismissed() {
		if opt.OnClick != nil && len(opt.Buttons) > 0 {
			opt.OnClick(0)
		}
	}
}
