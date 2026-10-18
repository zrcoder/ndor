package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

var ndor = New()

func main() {
	log.SetFlags(log.Lshortfile)

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Ndor",
			Width:   1200,
			Height:  800,
			Content: ui.View(ndor.View),
		})
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
