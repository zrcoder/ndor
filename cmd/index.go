package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/zrcoder/ndor/examples"
)

//go:embed *.gohtml
var indexFs embed.FS

var indexTmp = template.Must(template.ParseFS(indexFs, "index.gohtml"))

func index(w http.ResponseWriter, r *http.Request) {
	err := indexTmp.Execute(w, examples.Default)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
	}
}
