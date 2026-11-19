package comp

import (
	"github.com/egoist/mygo/ui"
)

type IconButtonOption struct {
	Icon       *ui.SVG
	SecondIcon *ui.SVG
	ShowSecond bool
	OnClick    func()
}

func IconButton(c *ui.Context, opt IconButtonOption) {
	btn := ui.ButtonBase(c).Margin(8).Padding(8).Radius(4).OnClick(func() {
		opt.OnClick()
	})
	if btn.Hovered() {
		btn.Background(c.Theme().SurfaceHover)
	}
	showIcon := opt.Icon
	if opt.SecondIcon != nil && opt.ShowSecond {
		showIcon = opt.SecondIcon
	}
	btn.Children(func() {
		ui.Icon(c, showIcon).Size(16, 16).TextColor(c.Theme().Text)
	})
}
