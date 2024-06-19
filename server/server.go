package server

import (
	"alnoor/todo-go-htmx/store"
	"alnoor/todo-go-htmx/views"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Store  store.SqliteStore
	Router *chi.Mux
}

func NewTasksServer() Server {
	store := store.SqliteStore{Path: "./todo.db"}
	server := Server{}
	store.Open()
	store.Migrate()
	server.Store = store

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.NoCache)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", server.IndexHandler)
	r.Post("/tasks", server.PostTaskHandler)
	server.Router = r

	return server

}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.Store.GetTasks()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	views.Index(tasks).Render(r.Context(), w)
}
func (s *Server) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("%v\n", err)
	}
	description := r.Form.Get("description")

	if description == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	task, err := s.Store.InsertTask(description)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("%v\n", err)
	}
	w.WriteHeader(http.StatusCreated)
	views.Task(task.Description).Render(r.Context(), w)
}
