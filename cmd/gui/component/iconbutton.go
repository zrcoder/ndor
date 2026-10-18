package component

import (
	"github.com/egoist/mygo/ui"
)

func IconButton(c *ui.Context, icon, checkIcon *ui.SVG, check bool, onClick func()) {
	btn := ui.ButtonBase(c).Margin(8).Padding(8).Radius(4).OnClick(func() {
		onClick()
	})
	if btn.Hovered() {
		btn.Background(ui.RGB(240, 240, 240))
	}
	showIcon := icon
	if checkIcon != nil && check {
		showIcon = checkIcon
	}
	btn.Children(func() {
		ui.Icon(c, showIcon).Size(16, 16).TextColor(ui.RGB(100, 100, 100))
	})
}
