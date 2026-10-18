package main

import (
	"embed"
	_ "embed"
	"os"
	"time"

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
	resultPNGData []byte
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
	copyCodeDone  time.Time
	copyIcon      *ui.SVG
	checkIcon     *ui.SVG
	saveIcon      *ui.SVG
	closeIcon     *ui.SVG
	sunIcon       *ui.SVG
	moonIcon      *ui.SVG
}

//go:embed resources/teacher.png
var teacherImgData []byte

//go:embed resources/teacher_labeled.png
var teacherLabeledImgData []byte

//go:embed resources/gopher_labeled.png
var gopherLabeledImgData []byte

//go:embed resources/wiki
var wikiFS embed.FS

//go:embed resources/icons/copy.svg
var copyIconData []byte

//go:embed resources/icons/check.svg
var checkIconData []byte

//go:embed resources/icons/save.svg
var saveIconData []byte

//go:embed resources/icons/close.svg
var closeIconData []byte

//go:embed resources/icons/sun.svg
var sunIconData []byte

//go:embed resources/icons/moon.svg
var moonIconData []byte

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
	n.copyIcon = ui.MustParseSVG(copyIconData)
	n.checkIcon = ui.MustParseSVG(checkIconData)
	n.saveIcon = ui.MustParseSVG(saveIconData)
	n.closeIcon = ui.MustParseSVG(closeIconData)
	n.sunIcon = ui.MustParseSVG(sunIconData)
	n.moonIcon = ui.MustParseSVG(moonIconData)
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
	// Handle keyboard shortcuts
	if c.Shortcut(ui.Ctrl, ui.KeyH) {
		n.showHelp = true
		c.Invalidate()
	}
	if c.Shortcut(ui.Ctrl, ui.KeyR) {
		n.runCode(c)
	}

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
			ui.Row(c).Fill().Height(40).Children(func() {
				ui.Spacer(c)
				component.IconButton(c, n.saveIcon, nil, false, func() {
					n.saveImage(c)
				})
			})
			ui.Image(c, n.resultImg).MaxWidthPercent(100).MaxHeightPercent(100)
		}
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
	ui.Column(c).Fill().BorderWidth(0, 0, 0, 0.5).BorderColor(component.DividerColor).Children(func() {
		ui.Row(c).Fill().Height(36).Children(func() {
			ui.Spacer(c)
			// Copy button
			copied := !n.copyCodeDone.IsZero() && c.Now().Sub(n.copyCodeDone) < 2*time.Second
			component.IconButton(c, n.copyIcon, n.checkIcon, copied, func() {
				c.Services().WriteClipboard(n.editor.Buffer.String())
				n.copyCodeDone = c.Now()
				c.Invalidate()
			})
			// Theme toggle button
			themeIcon := n.moonIcon
			isDarkMode := mygo.Theme.IsDark()
			if isDarkMode {
				themeIcon = n.sunIcon
			}
			component.IconButton(c, themeIcon, nil, false, func() {
				if isDarkMode {
					mygo.Theme.SetSource(mygo.ThemeLight)
				} else {
					mygo.Theme.SetSource(mygo.ThemeDark)
				}
				c.Invalidate()
			})
		})
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
		n.resultPNGData = nil
		c.Invalidate()
		return
	}

	n.errorMsg = ""
	n.errorLine = 0
	n.resultPNGData = pngData

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
	names := make([]string, len(examples.Default))
	for i, ex := range examples.Default {
		names[i] = ex.Name
	}

	cfg := component.OverlayPanelConfig{
		Title:   "Examples",
		Width:   180,
		OnClose: func() { n.showExamples = false },
	}
	component.OverlayPanel(c, cfg, n.closeIcon, names, func(item string) {
		for _, ex := range examples.Default {
			if ex.Name == item {
				n.code = ex.Code
				n.editor = &editor.State{Buffer: editor.NewBuffer(n.code)}
				n.editor.Buffer.SetLanguage(editor.LanguageNamed("c"))
				break
			}
		}
	})
}

func (n *Ndor) showHelpDocument(c *ui.Context) {
	content := loadDoc(n.isEnglish, n.helpDocTab)

	if ui.Modal(c, &n.showHelpDoc, func() {
		ui.Column(c).Width(800).Height(600).Radius(12).Background(c.Theme().Surface).Children(func() {
			// Tabs with language toggle and close button
			ui.Row(c).Fill().Height(40).Children(func() {
				ui.Tabs(c, &n.helpDocTab, "Home", "API", "Tutorials")
				ui.Spacer(c)
				// Language toggle button
				btn := ui.ButtonBase(c).Width(60).Height(28).Radius(4).
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
				if btn.Hovered() {
					btn.Background(c.Theme().SurfaceHover)
				}
				// Close button
				component.IconButton(c, n.closeIcon, nil, false, func() {
					n.showHelpDoc = false
					c.Invalidate()
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
