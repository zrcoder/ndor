package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var (
	Handler = &app.Handler{
		Name:        "Niudour",
		Description: "牛豆儿画图",
		Lang:        "zh_CN",
		Title:       "牛豆儿画图",
		Icon: app.Icon{
			Default: "images/gopher.png",
		},
		RawHeaders: []string{`
		<style>
			.editorLineErr {
				background: red;
			}
		</style>
		`},
		Styles: []string{"style.css"},
		Scripts: []string{
			"js/sweetalert2.js",
			"js/alert.js",
			"https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.35.0/min/vs/loader.min.js",
			"js/editor.js",
		},
	}
)
