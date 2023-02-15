build:
	go generate
	GOARCH=wasm GOOS=js go build -o static/web/app.wasm
	go build -o static/app

local: build
	cd static && ./app
	rm static/app

repo: build
	cd static && ./app niudour
	rm static/app

github: repo
	git add -A && \
	git commit -m '$(msg)' && \
	git push

netlify: local
	cd static && \
	git init && \
	git add -A && \
	git commit -m '$(msg)' && \
	git remote add origin https://github.com/zrcoder/niudour && \
	git push -f origin HEAD:netlify && \
	rm -rf .git

runcur:
	cd static && go run ../tools/server

run: local runcur

clear: 
	rm -rf static/web
	rm static/*.js
	rm static/*.html
	rm static/manifest*
	rm static/app.css