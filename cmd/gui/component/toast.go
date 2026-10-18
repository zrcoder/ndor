package component

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// Toast displays a centered toast message, auto-dismisses after delay seconds
func Toast(c *ui.Context, message string, delay int) {
	showToast(c, "", message, delay)
}

// ToastWithTitle displays a toast with a title
func ToastWithTitle(c *ui.Context, title, message string, delay int) {
	showToast(c, title, message, delay)
}

func showToast(c *ui.Context, title, message string, delay int) {
	if delay <= 0 {
		delay = ToastDelay
	}

	fullMessage := message
	if title != "" {
		fullMessage = title + ": " + message
	}

	ui.Row(c).Top(0).Bottom(0).Fill().Absolute().Children(func() {
		ui.Spacer(c)
		ui.Column(c).Width(ToastWidth).Height(ToastHeight).Radius(ToastRadius).Background(ErrorBgColor).Children(func() {
			ui.Row(c).Fill().Height(ToastHeight).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				ui.Text(c, fullMessage).TextColor(ErrorTextColor).FontSize(14)
			})
		})
		ui.Spacer(c)
	})

	// Auto dismiss
	go func() {
		time.Sleep(time.Duration(delay) * time.Second)
		c.Invalidate()
	}()
}

// ToastError displays an error toast
func ToastError(c *ui.Context, errMsg string, delay int) {
	if delay <= 0 {
		delay = ToastDelay
	}
	Toast(c, errMsg, delay)
}

// ToastErrorAtLine displays an error toast with line number
func ToastErrorAtLine(c *ui.Context, line int, errMsg string, delay int) {
	if delay <= 0 {
		delay = ToastDelay
	}
	msg := fmt.Sprintf("Error at line %d: %s", line, errMsg)
	Toast(c, msg, delay)
}
