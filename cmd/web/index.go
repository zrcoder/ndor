package main

import (
	"math/rand"
	"time"

	"github.com/zrcoder/ndor/examples"
	"github.com/zrcoder/ndor/internal"
	"github.com/zrcoder/ndor/pkg"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	pictureAreaID = "pictureArea"
	pictureBoxID  = "pictureBox"
)

const exampleCode = `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill

`

type index struct {
	app.Compo
	showExamples bool
}

func (idx *index) Render() app.UI {
	return app.Div().Style("overflow", "hidden").Body(
		app.Div().Class("title-bar").Body(
			app.P().Text("Ndor"),
			app.Button().Hidden(true).ID("example-button").OnClick(func(app.Context, app.Event) {
				idx.showExamples = true
			}),
		),
		app.Div().Class("main").Body(
			app.Div().ID(pictureBoxID).Class("left-box").Body(
				app.Img().ID(pictureAreaID).Style("max-width", "100%").Style("max-height", "100%"),
				app.If(idx.showExamples, func() app.UI {
					return app.Ul().Class("example-list").Body(
						app.Range(examples.Default).Slice(func(i int) app.UI {
							return app.Li().Text(examples.Default[i].Name).OnClick(func(_ app.Context, _ app.Event) {
								app.Window().Get("SetCode").Invoke(examples.Default[i].Code)
								idx.showExamples = false
							})
						}),
					)
				}),
			),
			app.Div().Class("right-box").ID("codeArea"),
		),
		app.Button().Class("teacher-button").OnClick(teacherButtonAction).Text("HELP"),
		app.Button().ID("run-button").Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func teacherButtonAction(ctx app.Context, e app.Event) {
	app.Window().Get("getndorAlert").Invoke().Call("showHelp")
}

func goButtonAction(ctx app.Context, e app.Event) {
	// TODO: disable gopher button when running draw proccess
	root := app.Window()
	root.Get(pictureAreaID).Set("src", "")
	alert := root.Get("getndorAlert").Invoke()
	alert.Call("toastPainting")

	pictureBox := root.Get(pictureBoxID)
	width := pictureBox.Get("offsetWidth").Int()
	height := pictureBox.Get("offsetHeight").Int()
	code := root.Get("GetCode").Invoke().String()

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
