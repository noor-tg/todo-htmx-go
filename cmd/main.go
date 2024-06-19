package main

import (
	"alnoor/todo-go-htmx/server"
	"net/http"
)

func main() {
	serve := server.NewTasksServer()
	http.ListenAndServe("localhost:3000", serve.Router)
}
