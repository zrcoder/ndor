package comp

import (
	"github.com/egoist/mygo/ui"
)

// OverlayPanelOption holds configuration for an overlay panel.
type OverlayPanelOption struct {
	Title          string
	Width          float32
	OnClose        func()
	OnClickOutside func() // called when user clicks outside the panel
	ItemFn         func(item string)
	CloseIcon      *ui.SVG
}

// OverlayPanel renders an overlay panel with a title bar and close button.
// Items are rendered via itemFn, which receives a *ui.Context for each item.
func OverlayPanel(c *ui.Context, opt OverlayPanelOption, items []string) {
	// Click anywhere to close
	ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Background(ui.Color{}).Absolute().OnClick(func() {
		if opt.OnClickOutside != nil {
			opt.OnClickOutside()
		}
		if opt.OnClose != nil {
			opt.OnClose()
		}
		c.Invalidate()
	})

	// Panel
	ui.Column(c).Top(10).Left(10).Width(opt.Width).Background(c.Theme().Surface).Radius(8).BorderWidth(1).BorderColor(c.Theme().Border).Absolute().Children(func() {
		// Header
		ui.Row(c).Fill().Height(36).BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border).Children(func() {
			ui.Text(c, opt.Title).FontSize(14).FontWeight(700).Padding(8)
			ui.Spacer(c)
			// Close button
			IconButton(c, IconButtonOption{
				Icon: opt.CloseIcon,
				OnClick: func() {
					if opt.OnClose != nil {
						opt.OnClose()
					}
					c.Invalidate()
				},
			})
			ui.Spacer(c).Width(4)
		})

		// Items
		ui.Column(c).Fill().Children(func() {
			for _, item := range items {
				itemCopy := item
				ui.Row(c).Width(opt.Width-10).Height(28).BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border).OnClick(func() {
					opt.ItemFn(itemCopy)
					if opt.OnClose != nil {
						opt.OnClose()
					}
					c.Invalidate()
				}).Children(func() {
					ui.Text(c, item).FontSize(13).Padding(0, 8, 0, 8)
				})
			}
		})
	})
}
