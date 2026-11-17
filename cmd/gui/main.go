package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/zrcoder/ndor/cmd/gui/turtle"
)

var ndor = New()

func main() {
	log.SetFlags(log.Lshortfile)

	mygo.Theme.SetSource(mygo.ThemeLight)

	turtleApp, err := turtle.New()
	if err != nil {
		log.Fatal(err)
	}

	var ndorWin, turtleWin *mygo.Window
	// Only one window is visible at a time; switching hides the other.
	showOnly := func(win, other *mygo.Window) {
		other.Hide()
		win.Show()
	}
	mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu},
		{Role: mygo.RoleFileMenu},
		{Role: mygo.RoleEditMenu},
		{Role: mygo.RoleViewMenu},
		{Label: "Window", Submenu: []*mygo.MenuItem{
			{ID: "ndor", Label: "ndor", Type: mygo.MenuItemRadio, Checked: true,
				Click: func(*mygo.MenuItem, *mygo.Window) { showOnly(ndorWin, turtleWin) }},
			{ID: "turtle", Label: "turtle", Type: mygo.MenuItemRadio,
				Click: func(*mygo.MenuItem, *mygo.Window) { showOnly(turtleWin, ndorWin) }},
		}},
	}))

	mygo.App.WhenReady(func() {
		ndorWin = mygo.NewWindow(mygo.WindowOptions{
			Title:   "Ndor",
			Width:   1200,
			Height:  800,
			Content: ui.View(ndor.View),
		})
		turtleWin = mygo.NewWindow(mygo.WindowOptions{
			Title:   "turtle",
			Width:   1600,
			Height:  900,
			Content: ui.View(turtleApp.View),
			Hidden:  true, // start hidden; only ndor shows until the menu switches
		})
		turtleApp.SetWindow(turtleWin)
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
