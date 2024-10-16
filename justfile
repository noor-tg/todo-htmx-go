# start dev
live:
	templ generate
	tailwind --no-autoprefixer -i assets/style.css -o static/style.css
	go build -o tmp/todo-dev cmd/main.go
	./tmp/todo-dev

# build production
build:
	go build -o tmp/todo cmd/main.go

# run tests & watch files and rerun tests
test:
	gow -g=richgo -c test -count=1 -v ./...
