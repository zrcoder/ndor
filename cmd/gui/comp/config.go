package comp

import "github.com/egoist/mygo/ui"

// Common color configuration
var (
	ErrorBgColor   = ui.RGB(80, 40, 40)
	ErrorTextColor = ui.RGB(255, 200, 200)

	DividerColor = ui.RGB(180, 180, 180)

	ButtonWidth  = float32(64)
	ButtonHeight = float32(72)

	ToastWidth  = float32(500)
	ToastHeight = float32(40)
	ToastRadius = float32(8)
	ToastDelay  = 3 // seconds
)
