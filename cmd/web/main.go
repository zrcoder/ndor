package main

import (
	"log"
	"os"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	version = "1.0.10"
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
		"https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.15/codemirror.min.css",
		"https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.15/theme/monokai.min.css",
		"https://cdnjs.cloudflare.com/ajax/libs/limonte-sweetalert2/11.7.1/sweetalert2.min.css",
		"style.css",
	},
	Scripts: []string{
		"https://cdnjs.cloudflare.com/ajax/libs/limonte-sweetalert2/11.7.1/sweetalert2.min.js",
		"js/alert.js",
		"https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.15/codemirror.min.js",
		"https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.15/mode/clike/clike.min.js",
		"js/editor.js",
	},
	Version: version,
}
