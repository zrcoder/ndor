package comp

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// mixedSegments builds a segmented control where each item can have icon,
// text label, or both — a hybrid neither Segments nor IconSegments covers.
// Copied from CrychicUI's internal segmentControl (segments.go); can be
// removed once that function is exported upstream.
func mixedSegments(c *ui.Context, label string, items []choice.SegmentItem, selected *int) {
	t := c.Theme()
	parts := ui.SegmentedBase(c, selected, len(items))
	tid, r := parts.Track.ID(), t.Radius-t.Space(0.5)
	parts.Track.Label(label).AlignSelf(ui.Start).Gap(t.Space(0.5)).
		Padding(t.Space(0.5)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
	parts.Track.Children(func() {
		for i, it := range items {
			c.Key(it.Label)
			seg := parts.Segment(i).Label(it.Label).Tooltip(it.Label).Radius(r)
			seg.Children(func() {
				ui.Grid(c).Columns(1).Children(func() {
					if i == *selected {
						ui.Box(c.Key("thumb")).ColumnStart(1).RowStart(1).Radius(r).
							Background(t.Background).Shadow(0, 1, 2, 0, kit.Shade(0.12)).
							Transition(ui.ElementTransition{Animation: &ui.Snappy, Match: thumbMatch{tid}})
					}
					face := ui.Row(c).ColumnStart(1).RowStart(1).Gap(t.Space(1.5)).
						Padding(t.Space(1.5)).AlignItems(ui.Center)
					face.Children(func() {
						if it.Icon != "" {
							ui.Icon(c, icon.Must(it.Icon)).Size(16, 16)
						} else {
							ui.Text(c, it.Label).SingleLine()
						}
					})
				})
			})
		}
	})
	kit.Changed(parts.Track, sound.Select)
}

type thumbMatch struct{ id uint64 }
