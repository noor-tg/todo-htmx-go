package main

import (
	"alnoor/todo-go-htmx"
	"alnoor/todo-go-htmx/server"
	"log"
	"net/http"
)

const (
	port = ":443"
)

func main() {
	serve := server.NewTasksServer(todo.ProductionCfg)

	srvr := &http.Server{
		Addr:    port,
		Handler: serve.Router,
	}

	log.Println("Starting HTTPS server...")
	err := srvr.ListenAndServeTLS("./todo.local.pem", "./todo.local-key.pem")
	if err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
