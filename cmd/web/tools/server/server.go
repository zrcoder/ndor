package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("serve on http://localhost:9999")
	log.Fatalln(http.ListenAndServe(":9999", http.FileServer(http.Dir("."))))
}
