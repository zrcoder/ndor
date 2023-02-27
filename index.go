package main

import (
	"math/rand"
	"time"

	"github.com/zrcoder/niudour/internal"
	"github.com/zrcoder/niudour/internal/config"
	"github.com/zrcoder/niudour/pkg"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

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

	showExamples bool
}

func (idx *index) Render() app.UI {
	return app.Div().Style("overflow", "hidden").Body(
		app.Div().Class("title-bar").Body(
			app.P().Text("Niudour 牛豆儿画图"),
			app.Button().Class("example-button").OnClick(func(ctx app.Context, e app.Event) {
				idx.showExamples = !idx.showExamples
			}).Text("Examples"),
		),
		app.Div().ID(pictureBoxID).Class("left-box").Body(
			app.Img().ID(pictureAreaID).Style("max-width", "100%").Style("max-height", "100%").Style("border", "4"),
		),
		app.Div().Class("right-box").Body(app.Pre().ID("codeArea").Class("code-area")),
		app.If(idx.showExamples, app.Ul().Class("example-list").Body(
			app.Range(config.Default.Examples).Slice(func(i int) app.UI {
				return app.Li().Text(config.Default.Examples[i].Name).OnClick(func(ctx app.Context, e app.Event) {
					app.Window().Get("SetCode").Invoke(config.Default.Examples[i].Code)
					idx.showExamples = false
				})
			}),
		)),
		app.Button().Class("teacher-button").OnClick(teacherButtonAction).Text("HELP"),
		app.Button().Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func teacherButtonAction(ctx app.Context, e app.Event) {
	alert := app.Window().Get("getNiudourAlert").Invoke()
	alert.Call("showHelp", version)
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

	// the draw proccess is very fast, wait a while to show the painting toast~
	time.Sleep(time.Duration(600+rand.Intn(800)) * time.Millisecond)
	alert.Call("closePaintToast")
	root.Get(pictureAreaID).Set("src", src)
}
