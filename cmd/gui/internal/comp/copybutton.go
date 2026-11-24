package comp

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CopyButton copies code to the clipboard and shows a check mark briefly.
type CopyButton struct {
	getContent func(*ui.Context) string
	done       time.Time
}

// NewCopyButton creates a CopyButton.
func NewCopyButton(getContent func(*ui.Context) string) *CopyButton {
	return &CopyButton{getContent: getContent}
}

// View renders the copy button.
func (b *CopyButton) View(c *ui.Context) {
	copied := !b.done.IsZero() && c.Now().Sub(b.done) < 2*time.Second
	if copied {
		c.After(2*time.Second - c.Now().Sub(b.done))
	}
	IconButton(c, IconButtonOption{
		Icon:       icon.Must("copy"),
		SecondIcon: icon.Must("check"),
		ShowSecond: copied,
		OnClick: func() {
			sound.Play(sound.Tap)
			c.Services().WriteClipboard(b.getContent(c))
			b.done = c.Now()
			c.Invalidate()
		},
	})
}
