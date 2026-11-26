package comp

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/navigation"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/zrcoder/ndor/cmd/gui/internal/locales"
)

// iconOption is one chip of an icon choice.
type iconOption struct{ value, icon, label string }

// Settings is the shared settings dialog: sound effects, appearance and
// language.
type Settings struct {
	open     bool
	soundOn  bool
	langOpen bool
}

// NewSettings creates settings reflecting the current app state.
func NewSettings() *Settings { return &Settings{soundOn: sound.Enabled()} }

// Open shows the settings dialog.
func (s *Settings) Open() { s.open, s.langOpen = true, false }

// View renders the dialog and applies changes.
func (s *Settings) View(c *ui.Context) {
	if !s.open {
		return
	}

	m, t := locales.In(c.Locale()), c.Theme()
	theme := string(mygo.Theme.Source())
	// The chips say their own language names; the tags are mygo's, so any
	// language other than Chinese falls back to English.
	langs, tags := []string{m.System(), "English", "中文"}, []string{"", "en", "zh-Hans"}
	chosen := max(slices.Index(tags, mygo.App.ChosenLocale()), 0)
	theme0, chosen0 := theme, chosen

	overlay.Modal(c, m.Settings(), &s.open, overlay.ModalOptions{Close: m.Close()}, func() {
		ui.Column(c).Width(360).Gap(t.Space(2)).Children(func() {
			typography.Label(c, m.Sound(), typography.LabelOptions{})
			choice.IconToggle(c, m.Sound(), "volume-x", "volume-2", &s.soundOn)

			typography.Label(c, m.Language(), typography.LabelOptions{})
			navigation.LanguageMenu(c, m.Language(), langs, &chosen, &s.langOpen)

			typography.Label(c, m.Theme(), typography.LabelOptions{})
			ui.RadioGroup(c, func() {
				for _, o := range []iconOption{
					{"system", "monitor", m.System()},
					{"light", "sun", m.Light()},
					{"dark", "moon", m.Dark()},
				} {
					c.Key(o.value)
					chip := ui.RadioBase(c, &theme, o.value).Label(o.label).Tooltip(o.label).Padding(t.Space(1), t.Space(2)).Radius(999)
					bg, fg := t.Surface, t.Text
					if o.value == theme {
						bg, fg = t.Accent.Alpha(0.14), t.Accent
					} else if chip.Hovered() {
						bg = t.SurfaceHover
					}
					chip.Background(chip.AnimateColor("bg", bg, ui.Snappy))
					kit.Changed(chip, sound.Select)
					chip.Children(func() {
						ui.Icon(c, icon.Must(o.icon)).Size(16, 16).TextColor(fg)
					})
				}
			}).Label(m.Theme()).Row().Gap(t.Space(2))
		})
	})

	if theme != theme0 {
		mygo.Theme.SetSource(mygo.ThemeSource(theme))
	}
	if chosen != chosen0 {
		mygo.App.SetLocale(tags[chosen])
	}
	if s.soundOn != sound.Enabled() {
		sound.SetEnabled(s.soundOn)
		if s.soundOn {
			// The switch's own cue was muted while sound was off.
			sound.Play(sound.Toggle)
		}
	}
}
