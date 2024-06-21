package main

import (
	"log"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	version = "0.70"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app.Route("/", func() app.Composer { return &index{} })
	app.RunWhenOnBrowser()
	if len(os.Args) > 1 { // for github/gitee
		handler.Resources = app.GitHubPages(os.Args[1])
	}
	if err := app.GenerateStaticWebsite(".", handler); err != nil {
		log.Fatal(err)
	}
}

var handler = &app.Handler{
	Name:        "ndor",
	Description: "Ndor Draw",
	Lang:        "zh_CN",
	Title:       "Ndor",
	Icon: app.Icon{
		Default: "images/gopher.png",
	},
	Styles: []string{"style.css"},
	Scripts: []string{
		"js/lib/sweetalert2.min.js", "js/alert.js",
		"js/lib/monaco-editor/vs/loader.js", "js/editor.js",
	},
	Version: version,
}
