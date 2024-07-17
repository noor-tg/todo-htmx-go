package main

import (
	"alnoor/todo-go-htmx"
	server_echo "alnoor/todo-go-htmx/server/echo"
	"log"
	"net/http"
)

const (
	port = ":3000"
)

func main() {
	serve := server_echo.NewTasksServer(todo.ProductionCfg)

	srvr := &http.Server{
		Addr:    port,
		Handler: serve.Router,
	}

	log.Println("Starting HTTPS server...")
	err := srvr.ListenAndServeTLS("todo.local.pem", "todo.local-key.pem")
	if err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
