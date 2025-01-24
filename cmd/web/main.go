package main

import (
	"log"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	version = "1.0.2"
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
	Title:       "Ndor",
	Description: "Draw",
	Lang:        "zh_CN",
	Icon: app.Icon{
		Default: "images/hi.png",
		SVG:     "images/hi.png", // not svg now, just to prevent the go-app's default one.
	},
	Styles: []string{
		"https://cdn.jsdelivr.net/npm/codemirror@5.65.15/lib/codemirror.min.css",
		"https://cdn.jsdelivr.net/npm/codemirror@5.65.15/theme/monokai.min.css",
		"https://cdn.jsdelivr.net/npm/sweetalert2@11.7.1/dist/sweetalert2.min.css",
		"style.css",
	},
	Scripts: []string{
		"https://cdn.jsdelivr.net/npm/sweetalert2@11.7.1/dist/sweetalert2.min.js",
		"js/alert.js",
		"https://cdn.jsdelivr.net/npm/codemirror@5.65.15/lib/codemirror.min.js",
		"https://cdn.jsdelivr.net/npm/codemirror@5.65.15/mode/clike/clike.min.js",
		"js/editor.js",
	},
	Version: version,
}
