package comp

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/egoist/mygo/ui"
)

type DialogOption struct {
	Title   string
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

	modal := overlay.Modal(c, opt.Title, opt.Open, overlay.ModalOptions{}, func() {
		ui.Column(c).Fill().AlignItems(ui.Center).Children(func() {
			if opt.Img != nil {
				ui.Image(c, opt.Img).Width(80).Height(100).Margin(0, 18)
			}

			ui.Row(c).Children(func() {
				for i, text := range opt.Buttons {
					fn := func() {
						*opt.Open = false
						if opt.OnClick != nil {
							opt.OnClick(i)
						}
					}
					if i == len(opt.Buttons)-1 {
						button.Solid(c, text).Margin(0, 10).OnClick(fn)
					} else {
						button.Outline(c, text).Margin(0, 10).OnClick(fn)
					}
				}
			})
		})
	})

	if modal.Changed() && opt.OnClick != nil {
		opt.OnClick(0)
	}
}
