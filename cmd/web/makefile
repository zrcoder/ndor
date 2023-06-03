build:
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

run:
	cd static && go run ../tools/server

clear: 
	rm -rf static/web
	rm static/*.js
	rm static/*.html
	rm static/manifest*
	rm static/app.css
	rm -rf static/.git
