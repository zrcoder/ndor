package main

import (
	"log"

	"github.com/zrcoder/niudour/internal"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app.Route("/", &index{})
	app.RunWhenOnBrowser()
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
	Version: "0.1",
}

const (
	pictureAreaID = "pictureArea"
	pictureBoxID  = "pictureBox"
)

const exampleCode = `context 800, 800
color 0, 255, 0, 255
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
		app.Button().Class("teacher-button").Attr("onclick", "alertTeacherHelp()").Text("HELP"),
		app.Button().Class("run-button").OnClick(goButtonAction).Text("GO"),
	)
}

func goButtonAction(ctx app.Context, e app.Event) {
	ctx.JSSrc().Set("disabled", "disabled") // TODO
	defer ctx.JSSrc().Set("disabled", "")

	root := app.Window()
	root.Get(pictureAreaID).Set("src", "")

	root.Get("toastPainting").Invoke()
	pictureBox := root.Get(pictureBoxID)
	width := pictureBox.Get("offsetWidth").Int()
	height := pictureBox.Get("offsetHeight").Int()
	code := root.Get("getCode").Invoke().String()
	src, errLine, errInfo := internal.Draw(width, height, code)
	root.Get("closePaintToast").Invoke()

	if errInfo != "" {
		if errInfo == internal.ErrEmptyInput.Error() {
			root.Get("alertEmptyInputWith").Invoke(exampleCode)
			return
		}
		root.Get("alertError").Invoke(errLine, errInfo)
		return
	}

	root.Get(pictureAreaID).Set("src", src)
}
