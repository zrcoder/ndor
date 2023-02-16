package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/zrcoder/niudour/internal"
	"github.com/zrcoder/niudour/pkg"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

//go:generate igop export -outdir ./internal/exported ./api

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
	Name:        "Niudour",
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
	Version: "0.38",
}

const (
	pictureAreaID = "pictureArea"
	pictureBoxID  = "pictureBox"
)

const exampleCode = `context 800, 800
color 0, 255, 0
circle 400, 400, 300
fill

// click the gopher on bottom right to draw!`

type index struct {
	app.Compo
}

func (i index) Render() app.UI {
	return app.Div().Style("overflow", "hidden").Body(
		app.Div().Class("title-bar").Body(app.P().Text("Niudour 牛豆儿画图")),
		app.Div().ID(pictureBoxID).Class("left-box").Body(
			app.Img().ID(pictureAreaID).Style("max-width", "100%").Style("max-height", "100%").Style("border", "4"),
		),
		app.Div().Class("right-box").Body(app.Pre().ID("codeArea").Class("code-area")),
		app.Button().Class("teacher-button").Attr("onclick", "niudourAlert.showHelp()").Text("HELP"),
		app.Button().Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func goButtonAction(ctx app.Context, e app.Event) {
	// TODO: disable gopher button when running draw proccess
	root := app.Window()
	root.Get(pictureAreaID).Set("src", "")
	alert := root.Get("getNiudourAlert").Invoke()

	pictureBox := root.Get(pictureBoxID)
	width := pictureBox.Get("offsetWidth").Int()
	height := pictureBox.Get("offsetHeight").Int()
	code := root.Get("GetCode").Invoke().String()

	alert.Call("toastPainting")

	src, err := pkg.Run(width, height, code)
	if err != nil {
		alert.Call("closePaintToast")
		if err == internal.ErrEmptyInput {
			alert.Call("alertEmptyInputWith", exampleCode)
			return
		}
		alert.Call("alertError", err.Number, err.Msg)
		return
	}

	// the draw proccess is very fast, wait 1-2 s to show painting toast~
	time.Sleep(time.Duration(1000+rand.Intn(1000)) * time.Millisecond)
	alert.Call("closePaintToast")
	root.Get(pictureAreaID).Set("src", src)
}
