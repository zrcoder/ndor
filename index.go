package main

import (
	"log"

	"gitee.com/rdor/niudour/internal"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const (
	pictureAreaID = "pictureArea"
	pictureBoxID  = "pictureBox"
)

type Index struct {
	app.Compo
}

func (i Index) Render() app.UI {
	return app.Div().Style("overflow", "hidden").Body(
		app.Div().Class("title-bar").Body(
			app.P().Text("Niudour 牛豆儿画图"),
		),
		app.Div().ID(pictureBoxID).Class("left-box").Body(
			app.Img().ID(pictureAreaID).Style("max-width", "100%").Style("max-height", "100%").Style("border", "4"),
		),
		app.Div().Class("right-box").Body(
			app.Pre().ID("codeArea").Class("code-area"),
		),
		app.Button().Class("teacher-button").Attr("onclick", "teacherAction()").Text("HELP"),
		app.Button().Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func goButtonAction(ctx app.Context, e app.Event) {
	root := app.Window()
	root.Get(pictureAreaID).Set("src", "")

	text := root.Get("getCode").Invoke().String()
	log.Println("receive code:")
	log.Println(text)

	pictureBox := root.Get(pictureBoxID)
	width := pictureBox.Get("offsetWidth").Int()
	height := pictureBox.Get("offsetHeight").Int()

	src, err := internal.Draw(width, height, text)
	if err != nil {
		root.Get("alertError").Invoke(err.Error())
		return
	}
	root.Get("toastSuccess").Invoke(src)
}
