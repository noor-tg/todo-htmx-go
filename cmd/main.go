package main

import (
	"alnoor/todo-go-htmx/server"
	"net/http"
)

func main() {
	serve := server.NewTasksServer(server.Config{Cleanup: false, LogHttp: true})
	http.ListenAndServe(":3000", serve.Router)
}
