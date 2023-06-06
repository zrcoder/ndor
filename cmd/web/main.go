package main

import (
	"log"
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const (
	version = "0.56"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app.Route("/", &index{})
	app.RunWhenOnBrowser()
	if len(os.Args) > 1 { // for github/gitee
		handler.Resources = app.CustomProvider("", os.Args[1])
	}
	if err := app.GenerateStaticWebsite(".", handler); err != nil {
		log.Fatal(err)
	}
}

var handler = &app.Handler{
	Name:        "ndor",
	Description: "牛豆儿画图",
	Lang:        "zh_CN",
	Title:       "牛豆儿画图",
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
