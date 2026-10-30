package comp

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ToolbarConfig holds the configuration for the Toolbar component.
type ToolbarConfig struct {
	CopyIcon  *ui.SVG
	CheckIcon *ui.SVG
	SaveIcon  *ui.SVG
	CloseIcon *ui.SVG

	GetCode func() string
	OnCopy  func(string)
	OnSave  func()
}

// Toolbar owns the code area toolbar state and actions.
type Toolbar struct {
	cfg      ToolbarConfig
	copyDone time.Time
}

// NewToolbar creates a new Toolbar component.
func NewToolbar(cfg ToolbarConfig) *Toolbar {
	return &Toolbar{cfg: cfg}
}

// View renders the toolbar.
func (t *Toolbar) View(c *ui.Context) {
	ui.Row(c).Fill().Height(36).Children(func() {
		ui.Spacer(c)

		// Copy button
		copied := !t.copyDone.IsZero() && c.Now().Sub(t.copyDone) < 2*time.Second
		if copied {
			c.After(2*time.Second - c.Now().Sub(t.copyDone))
		}
		copyIcon := t.cfg.CopyIcon
		secondIcon := t.cfg.CheckIcon
		IconButton(c, IconButtonOption{
			Icon:       copyIcon,
			SecondIcon: secondIcon,
			ShowSecond: copied,
			OnClick: func() {
				code := ""
				if t.cfg.GetCode != nil {
					code = t.cfg.GetCode()
				}
				c.Services().WriteClipboard(code)
				if t.cfg.OnCopy != nil {
					t.cfg.OnCopy(code)
				}
				t.copyDone = c.Now()
				c.Invalidate()
			},
		})
	})
}

// SaveButton renders just the save button for the picture area.
func (t *Toolbar) SaveButton(c *ui.Context) {
	if t.cfg.SaveIcon == nil {
		return
	}
	IconButton(c, IconButtonOption{
		Icon: t.cfg.SaveIcon,
		OnClick: func() {
			if t.cfg.OnSave != nil {
				t.cfg.OnSave()
			}
		},
	})
}
