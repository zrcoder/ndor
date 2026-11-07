package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

var turtle = NewTurtle()

func main() {
	log.SetFlags(log.Lshortfile)

	mygo.Theme.SetSource(mygo.ThemeLight)

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Turtle",
			Width:   1200,
			Height:  800,
			Content: ui.View(turtle.View),
		})
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
