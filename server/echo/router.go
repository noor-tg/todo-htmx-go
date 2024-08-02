package server_echo

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/store"
	"crypto/tls"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Store     store.SqliteStore
	Router    *echo.Echo
	TLSConfig *tls.Config
}

func NewTasksServer(config todo.Config) Server {
	store, _ := store.New(config.DB, config.Cleanup)
	server := Server{}
	store.Migrate()
	server.Store = store

	e := echo.New()

	if config.LogHttp {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())

	// import embedded directory with http.FS
	assetHandler := http.FileServer(http.FS(todo.Static))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("", assetHandler)))

	e.GET("/", server.IndexHandler)
	e.POST("/tasks", server.PostTaskHandler)
	e.GET("/tasks/:id", server.GetTaskFormHandler)
	e.PUT("/tasks/toggle-status/:id", server.ToggleStatusOfTaskHandler)
	e.PUT("/tasks/:id", server.UpdateTaskHandler)
	e.DELETE("/tasks/:id", server.DeleteTaskHandler)

	server.Router = e

	return server
}
