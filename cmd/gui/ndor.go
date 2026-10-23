package main

import (
	"os"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/zrcoder/ndor/cmd/gui/comp"
	"github.com/zrcoder/ndor/pkg"
	"github.com/zrcoder/ndor/pkg/examples"
)

type Ndor struct {
	resultImg     *ui.Bitmap
	resultPNGData []byte
	errorMsg      string
	errorLine     int
	editor        *editor.State
	helpBtnImg    *ui.Bitmap
	runBtnImg     *ui.Bitmap
	copyCodeDone  time.Time
	copyIcon      *ui.SVG
	checkIcon     *ui.SVG
	saveIcon      *ui.SVG
	sunIcon       *ui.SVG
	moonIcon      *ui.SVG
	paintImg      *ui.Bitmap
	codeImg       *ui.Bitmap
	// help owns the help dialog, documentation viewer and examples overlay.
	help *comp.Help
}

func New() *Ndor {
	n := new(Ndor)
	if bmp, err := ui.DecodeBitmap(teacherLabeledImgData); err == nil {
		n.helpBtnImg = bmp
	}
	if bmp, err := ui.DecodeBitmap(gopherLabeledImgData); err == nil {
		n.runBtnImg = bmp
	}
	n.copyIcon = ui.MustParseSVG(copyIconData)
	n.checkIcon = ui.MustParseSVG(checkIconData)
	n.saveIcon = ui.MustParseSVG(saveIconData)
	n.sunIcon = ui.MustParseSVG(sunIconData)
	n.moonIcon = ui.MustParseSVG(moonIconData)
	if bmp, err := ui.DecodeBitmap(paintImgData); err == nil {
		n.paintImg = bmp
	}
	if bmp, err := ui.DecodeBitmap(codeImgData); err == nil {
		n.codeImg = bmp
	}
	code := `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
`
	n.editor = &editor.State{Buffer: editor.NewBuffer(code)}
	n.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))

	var dialogImg *ui.Bitmap
	if bmp, err := ui.DecodeBitmap(teacherImgData); err == nil {
		dialogImg = bmp
	}
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

func (n *Ndor) View(c *ui.Context) {
	// Handle keyboard shortcuts before building the toast viewport, so a run
	// started here shows the painting in this same frame.
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		n.help.Open()
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		n.runCode(c)
	}

	// The painting toast keeps the result hidden until the animation is over.
	painting := comp.PaintingViewport(c, n.paintImg, n.codeImg)

	if n.errorMsg != "" {
		comp.ToastErrorAtLine(c, n.errorLine, n.errorMsg, func() {
			n.errorMsg = ""
			n.errorLine = 0
		})
	}

	// Help dialog, documentation and examples
	n.help.View(c)

	// Main layout
	ui.Row(c).Key("main").Fill().Children(func() {
		n.pictureArea(c, painting)
		n.codeArea(c)
	})

	// Floating buttons
	ui.Column(c).Key("help-float").Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.helpBtnImg, func() {
			n.help.Open()
		})
	})

	ui.Column(c).Key("run-float").Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.runBtnImg, func() {
			n.runCode(c)
		})
	})
}

func (n *Ndor) pictureArea(c *ui.Context, painting bool) {
	ui.Column(c).WidthPercent(50).Fill().Children(func() {
		if painting || n.resultImg == nil {
			return
		}
		ui.Row(c).Fill().Height(40).Children(func() {
			ui.Spacer(c)
			comp.IconButton(c, comp.IconButtonOption{
				Icon: n.saveIcon,
				OnClick: func() {
					n.saveImage(c)
				},
			})
		})
		ui.Image(c, n.resultImg).MaxWidthPercent(100).MaxHeightPercent(100)
	})
}

func (n *Ndor) saveImage(c *ui.Context) {
	if n.resultPNGData == nil {
		return
	}
	data := n.resultPNGData
	go func() {
		path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
			Title:       "Save Image",
			DefaultPath: "ndor.png",
			Filters:     []mygo.FileFilter{{Name: "PNG Image", Extensions: []string{"png"}}},
		})
		if err != nil || path == "" {
			return
		}
		os.WriteFile(path, data, 0o644)
	}()
}

func (n *Ndor) codeArea(c *ui.Context) {
	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Fill().Height(36).Children(func() {
			ui.Spacer(c)
			// Copy button
			copied := !n.copyCodeDone.IsZero() && c.Now().Sub(n.copyCodeDone) < 2*time.Second
			comp.IconButton(c, comp.IconButtonOption{
				Icon:       n.copyIcon,
				SecondIcon: n.checkIcon,
				ShowSecond: copied,
				OnClick: func() {
					c.Services().WriteClipboard(n.editor.Buffer.String())
					n.copyCodeDone = c.Now()
					c.Invalidate()
				},
			})
			// Theme toggle button
			themeIcon := n.moonIcon
			if c.Theme().Dark {
				themeIcon = n.sunIcon
			}
			comp.IconButton(c, comp.IconButtonOption{
				Icon: themeIcon,
				OnClick: func() {
					if mygo.Theme.IsDark() {
						mygo.Theme.SetSource(mygo.ThemeLight)
					} else {
						mygo.Theme.SetSource(mygo.ThemeDark)
					}
					c.Invalidate()
				},
			})
		})
		editor.View(c, n.editor).Fill()
	})
}

func (n *Ndor) runCode(c *ui.Context) {
	n.resultImg = nil
	c.Invalidate()

	width, height := getPictureSize(c)
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}
	code := n.editor.Buffer.String()

	// The toast owns the overlay and hides itself when its time is over.
	comp.ShowPainting(c)
	c.Invalidate()

	pngData, lineErr := pkg.Gen(width, height, code)
	if lineErr != nil {
		n.errorMsg = lineErr.Msg
		n.errorLine = lineErr.Number
		n.resultImg = nil
		n.resultPNGData = nil
	} else {
		n.errorMsg = ""
		n.errorLine = 0
		n.resultPNGData = pngData
		bitmap, err := ui.DecodeBitmap(pngData)
		if err != nil {
			n.errorMsg = "Failed to decode image: " + err.Error()
			n.errorLine = 0
			n.resultImg = nil
		} else {
			n.resultImg = bitmap
		}
	}
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
