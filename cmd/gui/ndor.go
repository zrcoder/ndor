package main

import (
	"embed"
	_ "embed"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/zrcoder/ndor/cmd/gui/component"
	"github.com/zrcoder/ndor/pkg"
	"github.com/zrcoder/ndor/pkg/examples"
)

type Ndor struct {
	code          string
	resultImg     *ui.Bitmap
	showHelp      bool
	showHelpDoc   bool
	helpDocTab    int
	showExamples  bool
	isEnglish     bool
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

//go:embed resources/wiki
var wikiFS embed.FS

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
	n.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))
	return n
}

func (n *Ndor) View(c *ui.Context) {
	if n.errorMsg != "" {
		component.ToastErrorAtLine(c, n.errorLine, n.errorMsg, 3, func() {
			n.errorMsg = ""
			n.errorLine = 0
		})
	}

	if n.showHelp {
		component.ShowDialogWithImage(c, &n.showHelp, n.helpDialogImg, []string{"Examples", "Document"}, func(btnIndex int) {
			if btnIndex == 0 {
				n.showExamples = true
			} else {
				n.showHelpDoc = true
			}
		}, func() {
			// Click outside triggers Examples button behavior
			n.showExamples = true
		})
	}

	// Help document dialog
	if n.showHelpDoc {
		n.showHelpDocument(c)
	}

	// Examples list overlay
	if n.showExamples {
		n.showExamplesList(c)
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

func (n *Ndor) showExamplesList(c *ui.Context) {
	// Click anywhere to close
	ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Background(ui.Color{}).Absolute().OnClick(func() {
		n.showExamples = false
		c.Invalidate()
	})

	// Examples list - narrow floating panel at top-left
	ui.Column(c).Top(10).Left(10).Width(180).Background(ui.RGB(255, 255, 255)).Radius(8).BorderWidth(1).BorderColor(ui.RGB(200, 200, 200)).Absolute().Children(func() {
		// Header
		ui.Row(c).Fill().Height(36).BorderWidth(0, 0, 1, 0).BorderColor(ui.RGB(220, 220, 220)).Children(func() {
			ui.Text(c, "Examples").FontSize(14).FontWeight(700).Padding(8)
			ui.Spacer(c)
			ui.ButtonBase(c).Width(24).Height(24).Background(ui.Color{}).OnClick(func() {
				n.showExamples = false
				c.Invalidate()
			}).Children(func() {
				ui.Text(c, "x").FontSize(14).TextColor(ui.RGB(100, 100, 100))
			})
		})

		// Example list
		ui.Column(c).Fill().Children(func() {
			for _, ex := range examples.Default {
				ui.Row(c).Width(170).Height(28).BorderWidth(0, 0, 1, 0).BorderColor(ui.RGB(235, 235, 235)).OnClick(func() {
					n.code = ex.Code
					n.editor = &editor.State{Buffer: editor.NewBuffer(n.code)}
					n.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))
					n.showExamples = false
					c.Invalidate()
				}).Children(func() {
					ui.Text(c, ex.Name).FontSize(13).Padding(0, 8, 0, 8)
				})
			}
		})
	})
}

func (n *Ndor) showHelpDocument(c *ui.Context) {
	content := loadDoc(n.isEnglish, n.helpDocTab)

	if ui.Modal(c, &n.showHelpDoc, func() {
		ui.Column(c).Width(800).Height(600).Radius(12).Background(ui.RGB(255, 255, 255)).Children(func() {
			// Tabs with language toggle
			ui.Row(c).Fill().Height(40).Children(func() {
				ui.Tabs(c, &n.helpDocTab, "Home", "API", "Tutorials")
				ui.Spacer(c)
				ui.ButtonBase(c).Width(60).Height(28).Radius(4).
					Background(ui.RGB(240, 240, 240)).
					OnClick(func() {
						n.isEnglish = !n.isEnglish
						c.Invalidate()
					}).Children(func() {
					if n.isEnglish {
						ui.Text(c, "🌐 EN").FontSize(12)
					} else {
						ui.Text(c, "🌐 中文").FontSize(12)
					}
				})
			})
			// Markdown content - scrollable
			ui.Scroll(c).Fill().Children(func() {
				ui.Column(c).Width(760).Padding(10, 0, 0, 0).Children(func() {
					markdown.View(c, content, nil)
				})
			})
		})
	}).Dismissed() {
		n.showHelpDoc = false
	}
}

func getPictureSize(c *ui.Context) (int, int) {
	win := mygo.FocusedWindow()
	if win != nil {
		bounds := win.ContentBounds()
		return bounds.Width * 55 / 100, bounds.Height - 40 - 72
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
