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

	helpBtnImg, _ := ui.DecodeBitmap(teacherLabeledImgData)
	runBtnImg, _ := ui.DecodeBitmap(gopherLabeledImgData)
	dialogImg, _ := ui.DecodeBitmap(teacherImgData)
	turtleApp, err := turtle.New(turtle.Config{
		HelpButtonImg: helpBtnImg,
		RunButtonImg:  runBtnImg,
		HelpDialogImg: dialogImg,
		CloseIcon:     ui.MustParseSVG(closeIconData),
		LoadDoc:       loadTurtleDoc,
	})
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
			Width:   1400,
			Height:  800,
			Content: ui.View(ndor.View),
		})
		turtleWin = mygo.NewWindow(mygo.WindowOptions{
			Title:   "Turtle",
			Width:   1400,
			Height:  800,
			Content: ui.View(turtleApp.View),
			Hidden:  true, // start hidden; only ndor shows until the menu switches
		})
		turtleApp.SetWindow(turtleWin)
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadTurtleDoc(isEnglish bool, tab int) string {
	base := "resources/wiki/turtle/"
	if isEnglish {
		base += "en/"
	}

	var name string
	switch tab {
	case 0:
		name = "Home.md"
	case 1:
		name = "Api.md"
	case 2:
		name = "Tutorials.md"
	}

	data, _ := wikiFS.ReadFile(base + name)
	return string(data)
}
