package main

import (
	"alnoor/todo-go-htmx/views"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", templ.Handler(views.Index()).ServeHTTP)

	http.ListenAndServe("localhost:3000", r)
}
