package comp

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// themes lists the appearance chips: the mygo.ThemeSource and its icon.
var themes = []struct{ name, icon string }{
	{"system", "monitor"},
	{"light", "sun"},
	{"dark", "moon"},
}

// Settings is the shared settings dialog: sound effects and appearance.
type Settings struct {
	open    bool
	soundOn bool
}

// NewSettings creates settings reflecting the current app state.
func NewSettings() *Settings {
	return &Settings{soundOn: sound.Enabled()}
}

// Open shows the settings dialog.
func (s *Settings) Open() { s.open = true }

// View renders the dialog and applies changes.
func (s *Settings) View(c *ui.Context) {
	if !s.open {
		return
	}

	theme := string(mygo.Theme.Source())
	t := c.Theme()
	overlay.Modal(c, "Settings", &s.open, overlay.ModalOptions{}, func() {
		ui.Column(c).Width(360).Gap(t.Space(2)).Children(func() {
			typography.Label(c, "Sound", typography.LabelOptions{})
			choice.IconToggle(c, "Sound", "volume-x", "volume-2", &s.soundOn)
			typography.Label(c, "Theme", typography.LabelOptions{})
			ui.RadioGroup(c, func() {
				for _, tm := range themes {
					c.Key(tm.name)
					chip := ui.RadioBase(c, &theme, tm.name).Label(tm.name).Padding(t.Space(1), t.Space(2)).Radius(999)
					bg, fg := t.Surface, t.Text
					if tm.name == theme {
						bg, fg = t.Accent.Alpha(0.14), t.Accent
					} else if chip.Hovered() {
						bg = t.SurfaceHover
					}
					chip.Background(chip.AnimateColor("bg", bg, ui.Snappy))
					if kit.Changed(chip, sound.Select) {
						mygo.Theme.SetSource(mygo.ThemeSource(tm.name))
					}
					chip.Children(func() {
						ui.Icon(c, icon.Must(tm.icon)).Size(16, 16).TextColor(fg)
					})
				}
			}).Label("Theme").Row().Gap(t.Space(2))
		})
	})

	if s.soundOn != sound.Enabled() {
		sound.SetEnabled(s.soundOn)
		if s.soundOn {
			sound.Play(sound.Toggle)
		}
	}
}
