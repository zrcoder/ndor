package comp

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// Toast displays a centered toast message
func Toast(c *ui.Context, message string, onDismiss func()) {
	showToast(c, "", message, onDismiss)
}

// ToastWithTitle displays a toast with a title
func ToastWithTitle(c *ui.Context, title, message string, onDismiss func()) {
	showToast(c, title, message, onDismiss)
}

// ToastErrorAtLine displays an error toast with line number
func ToastErrorAtLine(c *ui.Context, line int, errMsg string, onDismiss func()) {
	msg := fmt.Sprintf("Error at line %d: %s", line, errMsg)
	Toast(c, msg, onDismiss)
}

func showToast(c *ui.Context, title, message string, onDismiss func()) {
	fullMessage := message
	if title != "" {
		fullMessage = title + ": " + message
	}

	ui.Row(c).Top(0).Bottom(0).Fill().Absolute().OnClick(func() {
		if onDismiss != nil {
			onDismiss()
		}
		c.Invalidate()
	}).Children(func() {
		ui.Spacer(c)
		ui.Column(c).Width(ToastWidth).Height(ToastHeight).Radius(ToastRadius).Background(ErrorBgColor).Children(func() {
			ui.Row(c).Fill().Height(ToastHeight).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Text(c, fullMessage).TextColor(ErrorTextColor).FontSize(14)
			})
		})
		ui.Spacer(c)
	})
}
