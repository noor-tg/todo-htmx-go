package server

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/store"
	"alnoor/todo-go-htmx/views"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Store  store.SqliteStore
	Router *chi.Mux
}

func NewTasksServer(config todo.Config) Server {
	store := store.New(config.DB)
	server := Server{}
	store.Open(config.Cleanup)
	store.Migrate()
	server.Store = store

	r := chi.NewRouter()
	if config.LogHttp {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.NoCache)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", server.IndexHandler)
	r.Post("/tasks", server.PostTaskHandler)
	r.Get("/tasks/{id:[0-9]+}", server.GetTaskFormHandler)
	r.Put("/tasks/toggle-status/{id:[0-9]+}", server.ToggleStatusOfTaskHandler)
	r.Put("/tasks/{id:[0-9]+}", server.UpdateTaskHandler)
	r.Delete("/tasks/{id:[0-9]+}", server.DeleteTaskHandler)

	server.Router = r

	return server
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	activeStatus := r.URL.Query().Get("status")
	description := r.URL.Query().Get("description")

	var tasks []todo.Task
	var err error

	if activeStatus == "الكل" {
		activeStatus = ""
	}

	// not status, empty search
	if len(activeStatus) == 0 && len(description) == 0 {
		tasks, err = s.Store.GetTasks(nil)
	}
	// status, search
	if len(activeStatus) > 0 && len(description) > 0 {
		tasks, err = s.Store.GetTasks(map[string]string{
			"description": description,
			"status":      activeStatus,
		})
	}
	// not status, search
	if len(activeStatus) == 0 && len(description) > 0 {
		tasks, err = s.Store.GetTasks(map[string]string{
			"description": description,
		})
	}
	// status, no search
	if len(activeStatus) > 0 && len(description) == 0 {
		tasks, err = s.Store.GetTasks(map[string]string{
			"status": activeStatus,
		})
	}

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		views.ServerError().Render(r.Context(), w)
		return
	}

	if target, ok := r.Header["Hx-Target"]; ok {
		if target[0] == "list" {
			views.Tasks(tasks).Render(r.Context(), w)
		} else {
			views.Index(activeStatus, tasks).Render(r.Context(), w)
		}
	} else {
		views.Index(activeStatus, tasks).Render(r.Context(), w)
	}

}

func (s *Server) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	description := templ.EscapeString(r.Form.Get("description"))
	if description == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		views.ServerError().Render(r.Context(), w)
		return
	}

	task, err := s.Store.InsertTask(description)
	if err != nil {
		log.Printf("%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	views.Task(task).Render(r.Context(), w)
	views.InputForm(true).Render(r.Context(), w)
}

func (s *Server) GetTaskFormHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	task, err := s.Store.GetTaskById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	views.UpdateForm(task).Render(r.Context(), w)
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	description := templ.EscapeString(r.Form.Get("description"))
	if description == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		views.NotFound().Render(r.Context(), w)
		return
	}

	task, err := s.Store.UpdateTask(id, description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	views.Task(task).Render(r.Context(), w)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	err = s.Store.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (s *Server) ToggleStatusOfTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	task, err := s.Store.ToggleTaskStatus(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		views.ServerError().Render(r.Context(), w)
		log.Printf("%v\n", err)
		return
	}

	views.Task(task).Render(r.Context(), w)
}
