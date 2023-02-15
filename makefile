build:
	GOARCH=wasm GOOS=js go build -o static/web/app.wasm
	go build -o static/app

gen: build
	cd static && ./app
	rm static/app

run: gen
	cd static && go run ../server

clear: 
	rm -rf static/web
	rm static/*.js
	rm static/*.html
	rm static/manifest*
	rm static/app.css