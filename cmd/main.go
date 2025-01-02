package main

import (
	"log"
	"net/http"

	"github.com/zrcoder/ndor/cmd/static"

	"github.com/zrcoder/amisgo"
	"github.com/zrcoder/amisgo/conf"
)

func main() {
	log.SetFlags(log.Lshortfile)
	app := amisgo.New(
		conf.WithTitle("Ndor"),
		conf.WithLang("zh_CN"),
		conf.WithTheme(conf.ThemeDark),
		conf.WithIcon("/static/images/hi.png"),
	).
		StaticFS("/static", http.FS(static.FS)).
		Mount("/", index)

	log.Println("Server started at http://localhost:9999")
	log.Fatal(app.Run(":9999"))
}
