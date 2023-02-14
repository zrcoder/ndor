package main

import (
	"gitee.com/rdor/niudour/internal"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

const (
	pictureAreaID = "pictureArea"
	pictureBoxID  = "pictureBox"
)

type Index struct {
	app.Compo
}

func (i Index) Render() app.UI {
	return app.Div().Style("overflow", "hidden").Body(
		app.Div().Class("title-bar").Body(app.P().Text("Niudour 牛豆儿画图")),
		app.Div().ID(pictureBoxID).Class("left-box").Body(
			app.Img().ID(pictureAreaID).Style("max-width", "100%").Style("max-height", "100%").Style("border", "4"),
		),
		app.Div().Class("right-box").Body(app.Pre().ID("codeArea").Class("code-area")),
		app.Button().Class("teacher-button").Attr("onclick", "teacherAction()").Text("HELP"),
		app.Button().Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func goButtonAction(ctx app.Context, e app.Event) {
	ctx.JSSrc().Set("disabled", "disabled") // TODO
	defer ctx.JSSrc().Set("disabled", "")

	root := app.Window()
	root.Get(pictureAreaID).Set("src", "")
	root.Get("toastPaint").Invoke()

	pictureBox := root.Get(pictureBoxID)
	width := pictureBox.Get("offsetWidth").Int()
	height := pictureBox.Get("offsetHeight").Int()
	code := root.Get("getCode").Invoke().String()
	src, errLine, errInfo := internal.Draw(width, height, code)
	if errInfo != "" {
		root.Get("closeToastPaint").Invoke()
		root.Get("alertError").Invoke(errLine, errInfo)
		return
	}
	root.Get("closeToastPaint").Invoke()
	root.Get("alertSuccess").Invoke(src)
}
