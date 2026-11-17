package turtle

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/zrcoder/ndor/cmd/gui/turtle/api"
)

const (
	cell               = 10                    // 网格间距（DIPs），原来 40，现在更细密（4 倍线条）
	cw                 = 560                   // 画布宽
	ch                 = 440                   // 画布高
	stepDelay          = 80 * time.Millisecond // 每帧/每微步的间隔
	turnDelay          = 150 * time.Millisecond
	stepPixelsPerFrame = 4 // forward/backward 每帧移动的固定距离（DIPs），总时间与距离成正比
)

// stepOp is one micro-operation produced by the user's xgo code.
type stepOp struct {
	kind  string // "step" | "turn" | "penUp" | "penDown" | "setColor" | "setWidth" | "clear" | "home"
	val   float32
	color [3]uint8
}

type App struct {
	win       *mygo.Window
	t         turtle
	cmd       string
	stepQueue []stepOp // 按帧执行的微操作队列（不用 sleep，不用 goroutine）
	stepIdx   int
	stepTimer time.Time // 下一帧应执行的时间点
}

// New builds a turtle app; call SetWindow once its mygo window exists.
func New() (*App, error) {
	a := &App{}
	a.t.clear()
	if err := initGopher(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) SetWindow(w *mygo.Window) { a.win = w }

// pushStep 把一个微操作压入队列。
func (a *App) pushStep(op stepOp) {
	a.stepQueue = append(a.stepQueue, op)
}

// runXGoCode 解析并运行用户输入的 xgo 代码，将其转成按帧执行的微操作队列。
// 立即返回，不阻塞主线程。实际执行由 draw 里的 DrawOver 每帧推进。
func (a *App) runXGoCode(code string) {
	a.stepQueue = nil
	a.stepIdx = 0
	a.stepTimer = time.Time{}

	// 重置 api 命令队列
	api.CmdQueue = nil

	if err := runCode(code); err != nil {
		// 执行出错，显示错误信息
		a.cmd = code + "\n// ERROR: " + err.Error()
		return
	}

	// 把 api.CmdQueue 转成 stepOp 队列
	for _, c := range api.CmdQueue {
		switch c.Kind {
		case "step":
			// 用户单位 × cell = 实际 DIPs
			dist := float32(c.Val) * cell
			n := int(math.Ceil(float64(dist) / float64(stepPixelsPerFrame)))
			if n < 1 {
				n = 1
			}
			step := dist / float32(n)
			for i := 0; i < n; i++ {
				a.pushStep(stepOp{kind: "step", val: step})
			}
		case "turn":
			a.pushStep(stepOp{kind: "turn", val: float32(c.Val)})
		case "penUp":
			a.pushStep(stepOp{kind: "penUp"})
		case "penDown":
			a.pushStep(stepOp{kind: "penDown"})
		case "setColor":
			a.pushStep(stepOp{kind: "setColor", color: c.Color})
		case "setWidth":
			a.pushStep(stepOp{kind: "setWidth", val: float32(c.Val)})
		case "clear":
			a.pushStep(stepOp{kind: "clear"})
		case "home":
			a.pushStep(stepOp{kind: "home"})
		}
	}

	if len(a.stepQueue) > 0 {
		a.stepTimer = time.Now().Add(-stepDelay)
		a.win.Invalidate()
	}
}

// tickStep 由 draw 里的 DrawOver 每帧调用：若到了时间点，执行队列中下一个 op。
func (a *App) tickStep(p *ui.Painter) {
	if a.stepIdx >= len(a.stepQueue) {
		return
	}
	if !a.stepTimer.IsZero() && p.Now().Sub(a.stepTimer) < stepDelay {
		// 还没到时间，请求下一帧继续等待，保持动画流畅
		p.AnimationFrame()
		return
	}

	op := a.stepQueue[a.stepIdx]
	a.stepIdx++
	a.stepTimer = p.Now()

	// 直接在绘图回调里改状态，不再走 win.Update
	switch op.kind {
	case "step":
		a.t.step(op.val)
	case "turn":
		a.t.turn(op.val)
	case "penUp":
		a.t.penUp()
	case "penDown":
		a.t.penDown()
	case "setColor":
		a.t.setColor(op.color[0], op.color[1], op.color[2])
	case "setWidth":
		a.t.setWidth(op.val)
	case "clear":
		a.t.clear()
	case "home":
		a.t.home()
	}

	if a.stepIdx < len(a.stepQueue) {
		p.AnimationFrame()
	}
}

func boolWord(b bool, yes, no string) string {
	if b {
		return yes
	}
	return no
}

func (a *App) View(c *ui.Context) {
	// 左右布局：左边画布，右边编辑器与控件；两栏 fillHeight 并平分水平宽度。
	ui.Row(c).Fill().Gap(24).Padding(24).Children(func() {
		// 左：画布，填满左栏。
		ui.Box(c).WidthPercent(50).FillHeight().Background(ui.RGB(255, 255, 255)).DrawOver(func(p *ui.Painter, r ui.Rect) {
			// 按帧推进动画队列（不用 sleep，不用 goroutine）。
			a.tickStep(p)
			a.draw(p, r)
		})

		// 右：示例菜单、运行按钮、代码编辑器，填满右栏。
		ui.Column(c).WidthPercent(50).FillHeight().Gap(12).Children(func() {
			ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
				ui.MenuButton(c, "示例", func(m *ui.Menu) {
					for _, ex := range examples {
						if m.Item(ex.name).ID(ex.code).Chosen() {
							a.cmd = ex.code
							a.win.Invalidate()
						}
					}
				}).Width(120)
				if ui.Button(c, "运行").Clicked() {
					a.win.Update(func() { a.t.clear(); a.win.Invalidate() })
					a.runXGoCode(a.cmd)
				}
			})

			ui.TextArea(c, &a.cmd).FillWidth().FillHeight()

			pos := fmt.Sprintf("(%.0f, %.0f)  ∠%.0f°  笔%s", a.t.x, a.t.y, math.Mod(float64(a.t.angle), 360), boolWord(a.t.down, "↓", "↑"))
			ui.Text(c, pos).FontSize(12).TextColor(ui.RGB(140, 140, 140))
		})
	})
}

func (a *App) draw(p *ui.Painter, r ui.Rect) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2

	// 坐标网格：参考 goplayspace，按距中心格数分三级深浅，
	// 穿过原点的那条竖线和横线最深，作为坐标轴。
	gridLight := ui.RGB(240, 240, 240)  // 普通网格
	gridMid := ui.RGB(220, 220, 220)    // 每 5 格
	gridCenter := ui.RGB(150, 150, 150) // 中心线（坐标轴）

	strokeGrid := func(x0, y0, x1, y1 float32, col ui.Color) {
		ln := ui.Path{}
		ln.MoveTo(x0, y0).LineTo(x1, y1)
		p.StrokePath(&ln, 1, col)
	}

	// 垂直网格线，k=0 为原点（Y 轴）。
	nx := int(cx / cell)
	for k := -nx; k <= nx; k++ {
		x := cx + float32(k)*cell
		if x < r.X || x > r.X+r.W {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(x, r.Y, x, r.Y+r.H, col)
	}

	// 水平网格线，k=0 为原点（X 轴）。
	ny := int(cy / cell)
	for k := -ny; k <= ny; k++ {
		y := cy + float32(k)*cell
		if y < r.Y || y > r.Y+r.H {
			continue
		}
		col := gridLight
		if k == 0 {
			col = gridCenter
		} else if k%5 == 0 {
			col = gridMid
		}
		strokeGrid(r.X, y, r.X+r.W, y, col)
	}

	sx := func(q pt) float32 { return cx + q.X }
	sy := func(q pt) float32 { return cy - q.Y }

	strokeSeg := func(s seg) {
		if len(s.pts) < 2 {
			return
		}
		ln := ui.Path{}
		ln.MoveTo(sx(s.pts[0]), sy(s.pts[0]))
		for _, q := range s.pts[1:] {
			ln.LineTo(sx(q), sy(q))
		}
		p.StrokePath(&ln, s.width, ui.RGB(s.color[0], s.color[1], s.color[2]))
	}
	for _, s := range a.t.segments {
		strokeSeg(s)
	}
	if a.t.cur != nil {
		strokeSeg(*a.t.cur)
	}

	// 海龟本体：goplayspace 的 gopher 图，中心在 (px,py)，尖朝向 angle。
	px, py := sx(pt{a.t.x, a.t.y}), sy(pt{a.t.x, a.t.y})
	drawGopher(p, px, py, a.t.angle)
}
