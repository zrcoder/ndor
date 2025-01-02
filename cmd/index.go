package main

import (
	"errors"

	"github.com/zrcoder/ndor/examples"
	"github.com/zrcoder/ndor/pkg"

	"github.com/zrcoder/amisgo/comp"
)

const exampleCode = `context 800, 800
color "lightgreen"
circle 400, 400, 300
fill

`
const version = "0.9.2"

var index = comp.Page().Title(
	comp.Image().Width("30px").Height("40px").InnerClassName("border-none").Src("/static/images/good_morning.png"),
).Toolbar(
	comp.Group().Mode("normal").Body(
		comp.Select().Mode("inline").Name("examples").Label("Examples").Options(examples.Default...).Value(examples.DefaultCode),
		comp.Button().ActionType("url").Label("Documentation v"+version).Url("https://gitee.com/rdor/ndor/wikis"),
	),
).Body(
	comp.Form().WrapWithPanel(false).Body(
		comp.Group().ClassName("items-center").Body(
			comp.Image().InnerClassName("border-none ").Name("picture").ThumbMode("contain").ImageMode("original"),
			comp.Wrapper().ClassName("items-center w-full").Body(
				comp.Editor().Name("code").Language("c").Value("${examples}").Size("xxl").Options(comp.Schema{"fontSize": "15"}),
				comp.Button().ClassName("text-black hover:text-white").Style(comp.Schema{
					"background-image":    `url("/static/images/gopher.png")`,
					"width":               "80px",
					"height":              "88px",
					"padding-top":         "55px",
					"border":              0,
					"background-size":     "100% 100%",
					"background-repeat":   "no-repeat",
					"background-position": "bottom",
					"background-color":    "transparent",
					"transition":          "background-size 0.3s ease",
				}).Label("Go").ActionType("submit").Primary(true).Transform("code", "picture", "Done", func(input any) (any, error) {
					code := input.(string)
					width := 800
					height := 800
					src, err := pkg.Run(width, height, code)
					if err != nil {
						return nil, errors.New(err.Msg)
					}
					return src, nil
				}),
			),
		),
	).Actions(),
)
