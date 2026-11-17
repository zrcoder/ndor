package turtle

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

//go:embed gopher_walk.svg
var gopherFS embed.FS

// goplayspace 的 gopher_walk.svg 已裁剪为仅第一帧（50×50）。
const gopherDrawSize = 30 // 在画布上以多大画 gopher

var gopherWalk *ui.Bitmap // 仅第一帧的静态 gopher

func initGopher() error {
	walkData, err := gopherFS.ReadFile("gopher_walk.svg")
	if err != nil {
		return err
	}

	walkSVG, err := ui.ParseSVG(walkData)
	if err != nil {
		return err
	}

	// 把 SVG 栅格化成 RGBA，再包成 Bitmap（只做一次，后续复用）。
	gopherWalk = ui.NewBitmap(walkSVG.Image(50, 50, nil))
	return nil
}

// drawGopher 在 (x, y) 处按 angle 旋转画出 gopher（仅第一帧静态图）。
func drawGopher(p *ui.Painter, x, y float32, angle float32) {
	if gopherWalk == nil {
		return
	}

	// ui.Rotation 是屏幕坐标系顺时针；turtle angle 是数学角度（Y 向上，0 指右）。
	// goplayspace 的 gopher 默认面向上方，所以屏幕旋转角度 = 90 - angle。
	screenAngle := 90 - angle

	// 绕海龟中心 (x, y) 旋转，把 gopher 以其中心为锚点画在 (x, y)。
	p.Transform(ui.Rotation(screenAngle), x, y, func() {
		dst := ui.Rect{X: x - gopherDrawSize/2, Y: y - gopherDrawSize/2, W: gopherDrawSize, H: gopherDrawSize}
		p.ImageRegion(gopherWalk, ui.Rect{}, dst, ui.ImageOptions{})
	})
}
