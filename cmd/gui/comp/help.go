package comp

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// HelpConfig is what the Help component shows and does.
type HelpConfig struct {
	// DialogImg is the image shown in the help dialog.
	DialogImg *ui.Bitmap
	// CloseIcon is the icon of the close buttons.
	CloseIcon *ui.SVG
	// ExampleNames lists the examples, in order.
	ExampleNames []string
	// Examples maps each example name to its code.
	Examples map[string]string
	// OnPickExample is called with the code of the chosen example.
	OnPickExample func(code string)
	// LoadDoc returns the markdown of a documentation tab, in the chosen
	// language.
	LoadDoc func(isEnglish bool, tab int) string
}

// Help owns the help dialog, the documentation viewer and the examples
// overlay, including their open, tab and language state.
type Help struct {
	cfg          HelpConfig
	showDialog   bool
	showDoc      bool
	docTab       int
	showExamples bool
	isEnglish    bool
}

// NewHelp creates a Help component that shows what cfg says.
func NewHelp(cfg HelpConfig) *Help {
	return &Help{cfg: cfg}
}

// Open shows the help dialog.
func (h *Help) Open() {
	h.showDialog = true
}

// View renders the help dialog and whatever it opens. It draws nothing while
// they are closed.
func (h *Help) View(c *ui.Context) {
	if h.showDialog {
		ShowDialogWithImage(c, DialogOption{
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
	content := ""
	if h.cfg.LoadDoc != nil {
		content = h.cfg.LoadDoc(h.isEnglish, h.docTab)
	}

	if ui.Modal(c, &h.showDoc, func() {
		ui.Column(c).Width(800).Height(600).Radius(12).Background(c.Theme().Surface).Children(func() {
			// Tabs with language toggle and close button
			ui.Row(c).Fill().Height(40).Children(func() {
				ui.Tabs(c, &h.docTab, "Home", "API", "Tutorials")
				ui.Spacer(c)
				// Language toggle button
				btn := ui.ButtonBase(c).Width(60).Height(28).Radius(4).
					OnClick(func() {
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
				// Close button
				IconButton(c, IconButtonOption{
					Icon: h.cfg.CloseIcon,
					OnClick: func() {
						h.showDoc = false
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
		h.showDoc = false
	}
}

func (h *Help) showExamplesList(c *ui.Context) {
	OverlayPanel(c, OverlayPanelOption{
		Title:   "Examples",
		Width:   180,
		OnClose: func() { h.showExamples = false },
		ItemFn: func(item string) {
			code, ok := h.cfg.Examples[item]
			if !ok || h.cfg.OnPickExample == nil {
				return
			}
			h.cfg.OnPickExample(code)
		},
		CloseIcon: h.cfg.CloseIcon,
	}, h.cfg.ExampleNames)
}
