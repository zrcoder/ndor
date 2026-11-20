package comp

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// CopyButtonConfig holds the configuration for the CopyButton component.
type CopyButtonConfig struct {
	CopyIcon  *ui.SVG
	CheckIcon *ui.SVG
	GetCode   func() string
}

// CopyButton copies code to the clipboard and shows a check mark briefly.
type CopyButton struct {
	cfg  CopyButtonConfig
	done time.Time
}

// NewCopyButton creates a CopyButton.
func NewCopyButton(cfg CopyButtonConfig) *CopyButton {
	return &CopyButton{cfg: cfg}
}

// View renders the copy button.
func (b *CopyButton) View(c *ui.Context) {
	copied := !b.done.IsZero() && c.Now().Sub(b.done) < 2*time.Second
	if copied {
		c.After(2*time.Second - c.Now().Sub(b.done))
	}
	IconButton(c, IconButtonOption{
		Icon:       b.cfg.CopyIcon,
		SecondIcon: b.cfg.CheckIcon,
		ShowSecond: copied,
		OnClick: func() {
			c.Services().WriteClipboard(b.cfg.GetCode())
			b.done = c.Now()
			c.Invalidate()
		},
	})
}
