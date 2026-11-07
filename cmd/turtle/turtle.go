// Package main is the turtle graphics app: a xgo snippet drives a turtle
// that draws step by step on a coordinate grid, replayed as an animation.
package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/zrcoder/ndor/cmd/gui/comp"
	"github.com/zrcoder/ndor/pkg"
	"github.com/zrcoder/ndor/pkg/examples"
)

type TurtleApp struct {
	editor     *editor.State
	plan       *pkg.TurtlePlan
	runId      uint64
	errorMsg   string
	errorLine  int
	helpBtnImg *ui.Bitmap
	runBtnImg  *ui.Bitmap
	help       *comp.Help
}

func NewTurtle() *TurtleApp {
	n := new(TurtleApp)
	n.helpBtnImg, _ = ui.DecodeBitmap(teacherLabeledImgData)
	n.runBtnImg, _ = ui.DecodeBitmap(gopherLabeledImgData)

	code := examples.Default["turtle"]
	n.editor = &editor.State{Buffer: editor.NewBuffer(code)}
	n.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))

	var dialogImg *ui.Bitmap
	dialogImg, _ = ui.DecodeBitmap(teacherImgData)
	n.help = comp.NewHelp(comp.HelpConfig{
		DialogImg:    dialogImg,
		CloseIcon:    ui.MustParseSVG(closeIconData),
		ExampleNames: examples.DefaultNames,
		Examples:     examples.Default,
		OnPickExample: func(code string) {
			n.editor.Buffer.Set(code)
		},
		LoadDoc: loadDoc,
	})
	return n
}

func (n *TurtleApp) View(c *ui.Context) {
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		n.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		n.runCode(c)
	}

	if n.errorMsg != "" {
		comp.ToastErrorAtLine(c, n.errorLine, n.errorMsg, func() {
			n.errorMsg = ""
			n.errorLine = 0
		})
	}

	n.help.View(c)

	ui.Row(c.Key("main")).Fill().Children(func() {
		n.codeArea(c)
		n.pictureArea(c)
	})

	ui.Column(c.Key("help-float")).Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.helpBtnImg, func() {
			n.help.Open()
		})
	})

	ui.Column(c.Key("run-float")).Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.runBtnImg, func() {
			n.runCode(c)
		})
	})
}

func (n *TurtleApp) pictureArea(c *ui.Context) {
	ui.Column(c).WidthPercent(50).Fill().Children(func() {
		if n.plan != nil {
			comp.Turtle(c, n.plan, n.runId)
		}
	})
}

func (n *TurtleApp) codeArea(c *ui.Context) {
	ui.Column(c).Fill().Children(func() {
		editor.View(c, n.editor).Fill()
	})
}

func (n *TurtleApp) runCode(c *ui.Context) {
	width, height := getPictureSize(c)
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}
	code := n.editor.Buffer.String()

	plan, lineErr := pkg.TurtleGen(width, height, code)
	if lineErr != nil {
		n.errorMsg = lineErr.Msg
		n.errorLine = lineErr.Number
		n.plan = nil
		c.Invalidate()
		return
	}

	n.errorMsg = ""
	n.errorLine = 0
	n.plan = plan
	n.runId++
	c.Invalidate()
}

func getPictureSize(c *ui.Context) (int, int) {
	win := mygo.FocusedWindow()
	if win != nil {
		bounds := win.ContentBounds()
		return bounds.Width / 2, bounds.Height - 40 - 72
	}
	return 800, 600
}

func loadDoc(isEnglish bool, tab int) string {
	var base string
	if isEnglish {
		base = "resources/wiki/en/"
	} else {
		base = "resources/wiki/"
	}

	var name string
	switch tab {
	case 0:
		name = "Home.md"
	case 1:
		name = "Api.md"
	case 2:
		name = "Tutorials.md"
	}

	data, _ := wikiFS.ReadFile(base + name)
	return string(data)
}
