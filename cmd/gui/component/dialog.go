package component

import (
	"github.com/egoist/mygo/ui"
)

// DialogContentFunc is the content builder function for DialogWithContent
type DialogContentFunc func(c *ui.Context)

// DialogConfig holds dialog configuration
type DialogConfig struct {
	Title     string
	Content   any
	Buttons   []string
	OnConfirm func() // confirm button callback
	OnClose   func() // close callback
}

// DialogWithContent shows a custom dialog centered on screen without obscuring original content
func DialogWithContent(c *ui.Context, open *bool, title string, contentFunc DialogContentFunc, buttons []string) {
	if !*open {
		return
	}

	// Click overlay to close dialog (transparent overlay, only responds to clicks)
	ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Background(ui.Color{}).Absolute().OnClick(func() {
		*open = false
		c.Invalidate()
	})

	// Dialog body - centered vertically and horizontally
	ui.Column(c).Top(0).Bottom(0).Left(0).Right(0).Absolute().Children(func() {
		ui.Spacer(c)
		ui.Row(c).Fill().Children(func() {
			ui.Spacer(c)
			// Dialog container
			ui.Column(c).Width(500).Radius(12).BorderWidth(1).BorderColor(ui.RGB(200, 200, 200)).Background(ui.RGB(255, 255, 255)).Padding(16).Children(func() {
				// Title
				ui.Text(c, title).FontSize(18).FontWeight(700)

				// Custom content
				if contentFunc != nil {
					contentFunc(c)
				}

				// Button area
				if len(buttons) > 0 {
					ui.Row(c).Fill().Padding(8, 16, 16, 16).Justify(ui.End).Children(func() {
						for i, btnText := range buttons {
							idx := i
							isConfirm := idx == len(buttons)-1
							bgColor := ui.RGB(240, 240, 240)
							textColor := ui.RGB(60, 60, 60)
							if isConfirm {
								bgColor = ui.RGB(0, 122, 255) // confirm button blue
								textColor = ui.RGB(255, 255, 255)
							}
							ui.ButtonBase(c).Width(80).Height(36).Radius(6).Background(bgColor).OnClick(func() {
								*open = false
								c.Invalidate()
							}).Children(func() {
								ui.Text(c, btnText).FontSize(14).TextColor(textColor)
							})
							if idx < len(buttons)-1 {
								ui.Spacer(c).Width(8)
							}
						}
					})
				}
			})
			ui.Spacer(c)
		})
		ui.Spacer(c)
	})
}

// ShowDialogWithImage shows a dialog with an image (similar to web help dialog)
func ShowDialogWithImage(c *ui.Context, open *bool, img *ui.Bitmap, title string, buttons []string, onClick func(btnIndex int)) {
	if !*open {
		return
	}

	// Click overlay to close
	ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Background(ui.Color{}).Absolute().OnClick(func() {
		*open = false
		c.Invalidate()
	})

	// Dialog centered vertically and horizontally
	ui.Column(c).Top(0).Bottom(0).Left(0).Right(0).Absolute().Children(func() {
		ui.Spacer(c)
		ui.Row(c).Fill().Children(func() {
			ui.Spacer(c)
			ui.Column(c).Width(350).Radius(12).BorderWidth(1).BorderColor(ui.RGB(200, 200, 200)).Background(ui.RGB(255, 255, 255)).Padding(16).Children(func() {
				// Image
				if img != nil {
					ui.Image(c, img).Width(80).Height(100)
				}

				// Title
				ui.Text(c, title).FontSize(18).FontWeight(700).Padding(12, 0, 0, 0)

				// Buttons
				if len(buttons) > 0 {
					ui.Row(c).Fill().Padding(8, 16, 16, 16).Children(func() {
						for i, btnText := range buttons {
							idx := i
							bgColor := ui.RGB(240, 240, 240)
							textColor := ui.RGB(60, 60, 60)
							if i == len(buttons)-1 {
								bgColor = ui.RGB(0, 122, 255) // last button blue
								textColor = ui.RGB(255, 255, 255)
							}
							ui.ButtonBase(c).Width(100).Height(36).Radius(6).Background(bgColor).OnClick(func() {
								*open = false
								if onClick != nil {
									onClick(idx)
								}
								c.Invalidate()
							}).Children(func() {
								ui.Text(c, btnText).FontSize(14).TextColor(textColor)
							})
							if i < len(buttons)-1 {
								ui.Spacer(c).Width(8)
							}
						}
					})
				}
			})
			ui.Spacer(c)
		})
		ui.Spacer(c)
	})
}

// ShowDialog shows a simple dialog
func ShowDialog(c *ui.Context, open *bool, title string, content any, buttons ...string) {
	if !*open {
		return
	}

	// Click overlay to close (transparent background)
	ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Background(ui.Color{}).Absolute().OnClick(func() {
		*open = false
		c.Invalidate()
	})

	// Dialog centered vertically and horizontally
	ui.Column(c).Top(0).Bottom(0).Left(0).Right(0).Absolute().Children(func() {
		ui.Spacer(c)
		ui.Row(c).Fill().Children(func() {
			ui.Spacer(c)
			ui.Column(c).Width(500).Radius(12).BorderWidth(1).BorderColor(ui.RGB(200, 200, 200)).Background(ui.RGB(255, 255, 255)).Padding(16).Children(func() {
				// Title
				ui.Text(c, title).FontSize(18).FontWeight(700)

				// Content
				switch v := content.(type) {
				case string:
					ui.Text(c, v).FontSize(14).TextColor(ui.RGB(60, 60, 60)).Padding(0, 16, 12, 16)
				case *ui.Bitmap:
					ui.Image(c, v).MaxWidth(100).Padding(0, 16, 12, 16)
				}

				// Buttons
				if len(buttons) > 0 {
					ui.Row(c).Fill().Padding(8, 16, 16, 16).Justify(ui.End).Children(func() {
						for i, btnText := range buttons {
							isConfirm := i == len(buttons)-1
							bgColor := ui.RGB(240, 240, 240)
							textColor := ui.RGB(60, 60, 60)
							if isConfirm {
								bgColor = ui.RGB(0, 122, 255)
								textColor = ui.RGB(255, 255, 255)
							}
							ui.ButtonBase(c).Width(80).Height(36).Radius(6).Background(bgColor).OnClick(func() {
								*open = false
								c.Invalidate()
							}).Children(func() {
								ui.Text(c, btnText).FontSize(14).TextColor(textColor)
							})
							if i < len(buttons)-1 {
								ui.Spacer(c).Width(8)
							}
						}
					})
				}
			})
			ui.Spacer(c)
		})
		ui.Spacer(c)
	})
}
