package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/image", draw)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	log.Println("serve on http://localhost")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
