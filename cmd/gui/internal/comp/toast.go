package comp

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

const (
	toastWidth  = float32(500)
	toastHeight = float32(40)
	toastRadius = float32(8)
)

var (
	errorBgColor   = ui.RGB(80, 40, 40)
	errorTextColor = ui.RGB(255, 200, 200)
)

// Toast displays a centered toast message
func Toast(c *ui.Context, message string, onDismiss func()) {
	ui.Row(c).Top(0).Bottom(0).Fill().Absolute().OnClick(func() {
		if onDismiss != nil {
			onDismiss()
		}
		c.Invalidate()
	}).Children(func() {
		ui.Spacer(c)
		ui.Column(c).Width(toastWidth).Height(toastHeight).Radius(toastRadius).Background(errorBgColor).Children(func() {
			ui.Row(c).Fill().Height(toastHeight).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Text(c, message).TextColor(errorTextColor).FontSize(14)
			})
		})
		ui.Spacer(c)
	})
}

// ToastErrorAtLine displays an error toast with line number
func ToastErrorAtLine(c *ui.Context, line int, errMsg string, onDismiss func()) {
	Toast(c, fmt.Sprintf("Error at line %d: %s", line, errMsg), onDismiss)
}
