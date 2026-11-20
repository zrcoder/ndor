package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	log.SetFlags(log.Lshortfile)

	mygo.Theme.SetSource(mygo.ThemeLight)

	mygo.App.SetMenu(mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu},
	}))

	app := NewApp()

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "Ndor",
			Width:     1400,
			Height:    800,
			Content:   ui.View(app.view),
			Frameless: true,
		})
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
