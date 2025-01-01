package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	version = "0.9.0"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app.Route("/", func() app.Composer { return &index{} })
	app.RunWhenOnBrowser()
	http.Handle("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server started at http://localhost:9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

var handler = &app.Handler{
	Title:       "Ndor",
	Description: "Draw",
	Lang:        "zh_CN",
	Icon: app.Icon{
		Default: "/static/images/hi.png",
		SVG:     "/static/images/hi.png", // not svg now, just to prevent the go-app's default one.
	},
	Styles: []string{"/static/style.css"},
	Scripts: []string{
		"/static/js/lib/sweetalert2.min.js", "/static/js/alert.js",
		"/static/js/lib/monaco-editor/vs/loader.js", "/static/js/editor.js",
	},
	Version: version,
}
