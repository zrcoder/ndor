package main

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &Index{})
	app.RunWhenOnBrowser()
	if err := app.GenerateStaticWebsite(".", Handler); err != nil {
		log.Fatal(err)
	}
}
