live/templ:
	templ generate -v

live/styles:
	tailwind --no-autoprefixer -i assets/style.css -o static/style.css

live/server:
	go build -o tmp/main main.go

live:
	make -j3 live/server live/styles live/templ
