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

// ShowDialogWithImage shows a dialog with an image (similar to web help dialog)
// buttons order: first button appears on left, last button appears on right (and is styled as primary)
// onClickOutside is called when user clicks outside the dialog (like clicking the cancel button)
func ShowDialogWithImage(c *ui.Context, opt DialogOption) {
	if !*opt.Open {
		return
	}

	// Check if modal was dismissed (clicked outside or Escape)
	if ui.Modal(c, opt.Open, func() {
		// Dialog content
		ui.Column(c).Width(250).Radius(12).Background(c.Theme().Surface).Children(func() {
			// Image
			if opt.Img != nil {
				ui.Row(c).Fill().Children(func() {
					ui.Spacer(c)
					ui.Image(c, opt.Img).Width(80).Height(100)
					ui.Spacer(c)
				})
			}

			// Buttons - centered with spacers
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
