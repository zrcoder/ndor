package comp

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// Notice shows a message in the middle of the window for a while.
type Notice struct {
	tone        feedback.Tone
	title, body string
	dur         time.Duration // how long it shows
	at          time.Time     // when it started
}

// Show shows the notice for dur.
func (n *Notice) Show(tone feedback.Tone, title, body string, dur time.Duration) {
	n.tone, n.title, n.body, n.dur, n.at = tone, title, body, dur, time.Now()
	sound.Play(noticeCue(tone))
	if win := mygo.FocusedWindow(); win != nil {
		win.Invalidate()
	}
}

// noticeCue returns the sound of a tone.
func noticeCue(tone feedback.Tone) sound.Cue {
	switch tone {
	case feedback.ToneSuccess:
		return sound.Success
	case feedback.ToneWarning:
		return sound.Warning
	case feedback.ToneDanger:
		return sound.Error
	}
	return sound.Notify
}

// View renders the notice while its time lasts.
func (n *Notice) View(c *ui.Context) {
	left := n.dur - c.Now().Sub(n.at)
	if left <= 0 {
		return
	}
	c.After(left)
	ui.Row(c.Key("notice")).Absolute().Fill().Center().PassThrough().Children(func() {
		feedback.Alert(c, n.tone, n.title, n.body).
			MaxWidth(c.Theme().Space(60)).
			Background(c.Theme().Background)
	})
}
