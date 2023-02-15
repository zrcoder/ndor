build:
	GOARCH=wasm GOOS=js go build -o static/web/app.wasm
	go build -o static/app

local: build
	cd static && ./app
	rm static/app

repo: build
	cd static && ./app niudour
	rm static/app

run:
	cd static && go run ../server

clear: 
	rm -rf static/web
	rm static/*.js
	rm static/*.html
	rm static/manifest*
	rm static/app.css