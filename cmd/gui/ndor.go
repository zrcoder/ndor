package main

import (
	_ "embed"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/zrcoder/ndor/cmd/gui/component"
	"github.com/zrcoder/ndor/pkg"
)

type Ndor struct {
	code          string
	resultImg     *ui.Bitmap
	showHelp      bool
	errorMsg      string
	errorLine     int
	editor        *editor.State
	helpBtnImg    *ui.Bitmap
	helpDialogImg *ui.Bitmap
	runBtnImg     *ui.Bitmap
}

//go:embed resources/teacher.png
var teacherImgData []byte

//go:embed resources/teacher_labeled.png
var teacherLabeledImgData []byte

//go:embed resources/gopher_labeled.png
var gopherLabeledImgData []byte

func New() *Ndor {
	n := new(Ndor)
	if bmp, err := ui.DecodeBitmap(teacherLabeledImgData); err == nil {
		n.helpBtnImg = bmp
	}
	if bmp, err := ui.DecodeBitmap(gopherLabeledImgData); err == nil {
		n.runBtnImg = bmp
	}
	if bmp, err := ui.DecodeBitmap(teacherImgData); err == nil {
		n.helpDialogImg = bmp
	}
	n.code = `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
`
	n.editor = &editor.State{Buffer: editor.NewBuffer(n.code)}
	return n
}

func (n *Ndor) View(c *ui.Context) {
	if n.errorMsg != "" {
		component.ToastErrorAtLine(c, n.errorLine, n.errorMsg, 3)
	}

	if n.showHelp {
		component.ShowDialogWithImage(c, &n.showHelp, n.helpBtnImg, "Ndor Help", []string{"Document", "Examples"}, func(btnIndex int) {
			if btnIndex == 0 {
				// TODO: open document
			} else {
				// TODO: show examples
			}
		})
	}

	// Main layout
	ui.Row(c).Fill().Children(func() {
		n.pictureArea(c)
		n.codeArea(c)
	})

	// Floating buttons
	ui.Column(c).Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		component.ImageButtonSimple(c, n.helpBtnImg, func() {
			n.showHelp = true
		})
	})

	ui.Column(c).Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		component.ImageButtonSimple(c, n.runBtnImg, func() {
			n.code = n.editor.Buffer.String()
			n.runCode(c)
		})
	})
}

func (n *Ndor) pictureArea(c *ui.Context) {
	ui.Column(c).WidthPercent(55).Fill().Children(func() {
		if n.resultImg != nil {
			ui.Row(c).Fill().Children(func() {
				ui.Spacer(c)
				ui.Image(c, n.resultImg).MaxWidthPercent(100).MaxHeightPercent(100)
				ui.Spacer(c)
			})
		}
	})
}

func (n *Ndor) codeArea(c *ui.Context) {
	ui.Column(c).Fill().BorderWidth(0, 0, 0, 0.5).BorderColor(component.DividerColor).Children(func() {
		editor.View(c, n.editor).Fill()
	})
}

func (n *Ndor) runCode(c *ui.Context) {
	width, height := getPictureSize(c)
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}

	pngData, lineErr := pkg.Gen(width, height, n.code)
	if lineErr != nil {
		n.errorMsg = lineErr.Msg
		n.errorLine = lineErr.Number
		n.resultImg = nil
		c.Invalidate()
		return
	}

	n.errorMsg = ""
	n.errorLine = 0

	bitmap, err := ui.DecodeBitmap(pngData)
	if err != nil {
		n.errorMsg = "Failed to decode image: " + err.Error()
		n.errorLine = 0
		n.resultImg = nil
		c.Invalidate()
	} else {
		n.resultImg = bitmap
	}
	c.Invalidate()
}

func getPictureSize(c *ui.Context) (int, int) {
	win := mygo.FocusedWindow()
	if win != nil {
		bounds := win.ContentBounds()
		return bounds.Width * 55 / 100, bounds.Height - 40 - 72
	}
	return 800, 600
}
