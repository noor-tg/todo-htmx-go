package main

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"net/http"
)

func main() {
	serve := server.NewTasksServer(todo.ProductionCfg)
	http.ListenAndServe(":3000", serve.Router)
}
