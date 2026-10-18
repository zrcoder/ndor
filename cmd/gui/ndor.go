package main

import (
	"os"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/zrcoder/ndor/cmd/gui/comp"
	"github.com/zrcoder/ndor/pkg"
	"github.com/zrcoder/ndor/pkg/examples"
)

type Ndor struct {
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
	paintImg      *ui.Bitmap
	// runSeq increments when a run starts; resultSeq follows it when the
	// picture of that run is ready (or it failed). The run is in flight
	// while they differ.
	runSeq    uint64
	resultSeq uint64
}

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
	if bmp, err := ui.DecodeBitmap(paintImgData); err == nil {
		n.paintImg = bmp
	}
	code := `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill
`
	n.editor = &editor.State{Buffer: editor.NewBuffer(code)}
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
		comp.ToastErrorAtLine(c, n.errorLine, n.errorMsg, func() {
			n.errorMsg = ""
			n.errorLine = 0
		})
	}

	if n.showHelp {
		opt := comp.DialogOption{
			Open:    &n.showHelp,
			Img:     n.helpDialogImg,
			Buttons: []string{"Examples", "Document"},
			OnClick: func(btnIndex int) {
				if btnIndex == 0 {
					n.showExamples = true
				} else {
					n.showHelpDoc = true
				}
			},
		}
		comp.ShowDialogWithImage(c, opt)
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
	ui.Row(c).Key("main").Fill().Children(func() {
		n.pictureArea(c)
		n.codeArea(c)
	})

	// Floating buttons
	ui.Column(c).Key("help-float").Width(64).Height(72).Bottom(0).Left(4).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.helpBtnImg, func() {
			n.showHelp = true
		})
	})

	ui.Column(c).Key("run-float").Width(64).Height(72).Bottom(0).Right(0).Absolute().Children(func() {
		comp.ImageButtonSimple(c, n.runBtnImg, func() {
			n.runCode(c)
		})
	})
}

// paintDuration is how long the painting animation runs before the picture
// of the run replaces it.
const paintDuration = 900 * time.Millisecond

func (n *Ndor) pictureArea(c *ui.Context) {
	ui.Column(c).WidthPercent(55).Fill().Children(func() {
		// The painting animation: a timeline keyed on the run, so a new run
		// plays it again. It shows while the run is in flight and while the
		// timeline has not come to rest, so the picture updates only once
		// both the animation and the generation are done.
		overlay := ui.Column(c).Key("painting")
		tl := overlay.Timeline("painting", ui.TimelineOptions{Run: n.runSeq})
		inFlight := n.resultSeq != n.runSeq
		painting := n.runSeq > 0 && (inFlight || !tl.Settled())
		if painting {
			if inFlight {
				// Keep frames coming until the result arrives; the timeline
				// only asks for frames while it moves.
				c.AnimationFrame()
			}
			progress := tl.Tween(ui.Tween{From: 0, To: 1, Duration: paintDuration, Ease: ui.Linear})
			overlay.Fill().Center().Gap(16).Children(func() {
				if n.paintImg != nil {
					ui.Image(c, n.paintImg).MaxHeightPercent(40)
				}

				ui.Progress(c, float64(progress)).Width(200)
			})
			return
		}

		if n.resultImg == nil {
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
	ui.Column(c).Fill().BorderWidth(0, 0, 0, 0.5).BorderColor(comp.DividerColor).Children(func() {
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
	if n.resultSeq != n.runSeq {
		return // a run is already in flight
	}

	width, height := getPictureSize(c)
	if width <= 0 || height <= 0 {
		width, height = 800, 800
	}
	code := n.editor.Buffer.String()

	// Bump the sequence: the view starts the painting animation for this
	// run, and the goroutine reports back under the same number.
	n.runSeq++
	seq := n.runSeq
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
	n.resultSeq = seq
}

func (n *Ndor) showExamplesList(c *ui.Context) {
	opt := comp.OverlayPanelOption{
		Title:   "Examples",
		Width:   180,
		OnClose: func() { n.showExamples = false },
		ItemFn: func(item string) {
			code, ok := examples.Default[item]
			if ok {
				n.editor.Buffer.Set(code)
			}
		},
		CloseIcon: n.closeIcon,
	}
	comp.OverlayPanel(c, opt, examples.DefaultNames)
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
				comp.IconButton(c, comp.IconButtonOption{
					Icon: n.closeIcon,
					OnClick: func() {
						n.showHelpDoc = false
						c.Invalidate()
					},
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
