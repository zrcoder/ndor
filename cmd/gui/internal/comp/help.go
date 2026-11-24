package comp

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/navigation"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// HelpConfig configures the Help component.
type HelpConfig struct {
	DialogImg     *ui.Bitmap
	ExampleNames  []string
	Examples      map[string]string
	OnPickExample func(code string)
	LoadDoc       func(isEnglish bool, tab int) string
}

// Help owns the help dialog, docs viewer and examples overlay state.
type Help struct {
	cfg          HelpConfig
	showDialog   bool
	showDoc      bool
	docTab       int
	showExamples bool
	isEnglish    bool
	exampleList  ui.ListState
}

// NewHelp creates a Help component.
func NewHelp(cfg HelpConfig) *Help {
	return &Help{cfg: cfg}
}

// Open shows the help dialog.
func (h *Help) Open() {
	h.showDialog = true
}

// View renders the help dialog and whatever it opens.
func (h *Help) View(c *ui.Context) {
	if h.showDialog {
		ShowDialogWithImage(c, DialogOption{
			Title:   "Help",
			Open:    &h.showDialog,
			Img:     h.cfg.DialogImg,
			Buttons: []string{"Examples", "Document"},
			OnClick: func(btn int) {
				if btn == 0 {
					h.showExamples = true
				} else {
					h.showDoc = true
				}
			},
		})
	}

	if h.showDoc {
		h.showDocument(c)
	}

	if h.showExamples {
		h.showExamplesList(c)
	}
}

func (h *Help) showDocument(c *ui.Context) {
	doc := ""
	if h.cfg.LoadDoc != nil {
		doc = h.cfg.LoadDoc(h.isEnglish, h.docTab)
	}

	closeFn := func() {
		sound.Play(sound.Close)
		h.showDoc = false
	}

	if ui.Modal(c, &h.showDoc, func() {
		ui.Column(c).Width(800).Height(600).Radius(12).Background(c.Theme().Surface).Children(func() {
			ui.Row(c).Fill().Height(40).Children(func() {
				navigation.UnderlineTabs(c, "pages", []string{"Home", "API", "Tutorials"}, &h.docTab)
				ui.Spacer(c)
				btn := ui.ButtonBase(c).Width(60).Height(28).Radius(4).
					OnClick(func() {
						sound.Play(sound.Tap)
						h.isEnglish = !h.isEnglish
						c.Invalidate()
					}).Children(func() {
					if h.isEnglish {
						ui.Text(c, "🌐 EN").FontSize(12)
					} else {
						ui.Text(c, "🌐 中文").FontSize(12)
					}
				})
				if btn.Hovered() {
					btn.Background(c.Theme().SurfaceHover)
				}
				IconButton(c, IconButtonOption{
					Icon:    icon.Must("x"),
					OnClick: closeFn,
				})
			})
			layout.ScrollArea(c, func() {
				ui.Column(c).Width(760).Padding(10, 0, 0, 0).Children(func() {
					markdown.View(c, doc, &markdown.Options{})
				})
			}).Fill()
		})
	}).Dismissed() {
		closeFn()
	}
}

func (h *Help) showExamplesList(c *ui.Context) {
	names := h.cfg.ExampleNames
	closePanel := func() {
		h.showExamples = false
		c.Invalidate()
	}

	ui.Overlay(c, func() {
		ui.Row(c).Top(0).Bottom(0).Left(0).Right(0).Absolute().OnClick(closePanel)

		ui.Column(c).Top(52).Left(10).Width(180).Background(c.Theme().Surface).Radius(8).BorderWidth(1).BorderColor(c.Theme().Border).Absolute().Children(func() {
			ui.Row(c).Fill().Height(36).BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border).Children(func() {
				ui.Text(c, "Examples").FontSize(14).FontWeight(700).Padding(8)
				ui.Spacer(c)
				IconButton(c, IconButtonOption{Icon: icon.Must("x"), OnClick: closePanel})
				ui.Spacer(c).Width(4)
			})

			ui.List(c, &h.exampleList, len(names), func(i int) {
				name := names[i]
				ui.Row(c).Width(170).Height(28).BorderWidth(0, 0, 1, 0).BorderColor(c.Theme().Border).OnClick(func() {
					sound.Play(sound.Select)
					if h.cfg.OnPickExample != nil {
						if code, ok := h.cfg.Examples[name]; ok {
							h.cfg.OnPickExample(code)
						}
					}
					closePanel()
				}).Children(func() {
					ui.Text(c, name).FontSize(13).Padding(0, 8, 0, 8)
				})
			}).Height(float32(len(names) * 28))
		})
	})
}
