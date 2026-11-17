package turtle

import "math"

// pt 是海龟的数学坐标：原点在画布中心，Y 向上。
type pt struct{ X, Y float32 }

// seg 是一条连续笔迹，带自己的颜色和宽度。
type seg struct {
	pts   []pt
	color [3]uint8
	width float32
}

type turtle struct {
	x, y, angle float32 // angle 度，0 指向右（东）
	down        bool
	penColor    [3]uint8
	penWidth    float32
	segments    []seg // 已完成的段
	cur         *seg  // 当前正在画的段（笔下时非 nil）
}

const defaultWidth = 2

func (t *turtle) turn(d float32) { t.angle += d }

func (t *turtle) penDown() {
	t.down = true
	if t.cur == nil {
		t.cur = &seg{color: t.penColor, width: t.penWidth}
		t.cur.pts = append(t.cur.pts, pt{t.x, t.y})
	}
}

func (t *turtle) penUp() {
	t.down = false
	if t.cur != nil && len(t.cur.pts) > 0 {
		t.segments = append(t.segments, *t.cur)
		t.cur = nil
	}
}

func (t *turtle) home() {
	t.x, t.y, t.angle = 0, 0, 0
	if t.down {
		t.penUp()
		t.penDown()
	}
}

func (t *turtle) clear() {
	t.x, t.y, t.angle = 0, 0, 0
	t.down = true
	t.penColor = [3]uint8{20, 20, 20}
	t.penWidth = defaultWidth
	t.segments = nil
	t.cur = nil
	t.penDown()
}

func (t *turtle) setColor(r, g, b uint8) {
	t.penColor = [3]uint8{r, g, b}
	if t.cur != nil && len(t.cur.pts) > 1 {
		// 不抬笔改颜色：先把当前段存入历史，再开新段，避免旧点染上新颜色
		t.segments = append(t.segments, *t.cur)
		t.cur = nil
	}
	if t.down && t.cur == nil {
		t.penDown()
	}
}

func (t *turtle) setWidth(w float32) {
	t.penWidth = w
	if t.cur != nil && len(t.cur.pts) > 1 {
		t.segments = append(t.segments, *t.cur)
		t.cur = nil
	}
	if t.down && t.cur == nil {
		t.penDown()
	}
}

// step 向前移动 d（负值即向后），笔下时把新点追加到当前段。
func (t *turtle) step(d float32) {
	rad := t.angle * math.Pi / 180
	nx := t.x + float32(math.Cos(float64(rad)))*d
	ny := t.y + float32(math.Sin(float64(rad)))*d
	if t.down {
		if t.cur == nil {
			t.penDown()
		}
		t.cur.pts = append(t.cur.pts, pt{nx, ny})
	}
	t.x, t.y = nx, ny
}
