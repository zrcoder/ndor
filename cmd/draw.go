package main

import (
	"encoding/json"
	"net/http"

	"github.com/zrcoder/ndor/internal"
	"github.com/zrcoder/ndor/pkg"
)

type Request struct {
	Width  int
	Height int
	Code   string
}

type Response struct {
	Error *internal.LineError
	Data  string
}

func draw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "you need post your code to draw", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer func() {
		r.Body.Close()
	}()
	req := &Request{}
	err := decoder.Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imageData, lineErr := pkg.Run(req.Width, req.Height, req.Code)
	resp := &Response{
		Error: lineErr,
		Data:  imageData,
	}
	data, _ := json.Marshal(resp)
	w.Write(data)
}
