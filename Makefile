live/templ:
	templ generate

live/styles:
	tailwind --no-autoprefixer -i assets/style.css -o static/style.css

live/server:
	go build -o tmp/main cmd/main.go

live:
	make -j3 live/server live/styles live/templ

build:
	go build -o tmp/todo cmd/main.go

test:
	gow -g=richgo -c test -count=1 -v ./...
