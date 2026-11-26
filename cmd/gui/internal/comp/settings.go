package comp

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/overlay"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/zrcoder/ndor/cmd/gui/internal/locales"
)

var (
	themeVals = []string{"system", "light", "dark"}
	langTags  = []string{"", "en", "zh-Hans"}
)

type Settings struct {
	open     bool
	soundOn  bool
	themeIdx int
	langIdx  int
}

func NewSettings() *Settings {
	return &Settings{
		soundOn:  sound.Enabled(),
		themeIdx: slices.Index(themeVals, string(mygo.Theme.Source())),
		langIdx:  slices.Index(langTags, mygo.App.ChosenLocale()),
	}
}

func (s *Settings) Open() { s.open = true }

func (s *Settings) View(c *ui.Context) {
	if !s.open {
		return
	}
	m, t := locales.In(c.Locale()), c.Theme()
	theme0, lang0 := s.themeIdx, s.langIdx

	overlay.Modal(c, m.Settings(), &s.open, overlay.ModalOptions{Close: m.Close()}, func() {
		ui.Grid(c).Columns(2).Gap(t.Space(4)).Children(func() {
			// Sound (spans both cols).
			ui.Column(c).ColumnSpan(2).Gap(t.Space(1)).Children(func() {
				typography.Label(c, m.Sound(), typography.LabelOptions{})
				ic := "volume-x"
				if s.soundOn {
					ic = "volume-2"
				}
				ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
					ui.Icon(c, icon.Must(ic)).Size(16, 16).TextColor(t.TextMuted)
					choice.Toggle(c, " ", &s.soundOn)
				})
			})

			// Theme.
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				typography.Label(c, m.Theme(), typography.LabelOptions{})
				choice.IconSegments(c, m.Theme(), []choice.SegmentItem{
					{Label: m.System(), Icon: "monitor"},
					{Label: m.Light(), Icon: "sun"},
					{Label: m.Dark(), Icon: "moon"},
				}, &s.themeIdx)
			})

			// Language — mixed: icon for System, text for En/文.
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				typography.Label(c, m.Language(), typography.LabelOptions{})
				mixedSegments(c, m.Language(), []choice.SegmentItem{
					{Label: m.System(), Icon: "monitor"},
					{Label: "En"}, {Label: "文"},
				}, &s.langIdx)
			})
		})
	})

	if s.themeIdx != theme0 {
		mygo.Theme.SetSource(mygo.ThemeSource(themeVals[s.themeIdx]))
	}
	if s.langIdx != lang0 {
		mygo.App.SetLocale(langTags[s.langIdx])
	}
	if s.soundOn != sound.Enabled() {
		sound.SetEnabled(s.soundOn)
		if s.soundOn {
			sound.Play(sound.Toggle)
		}
	}
}
