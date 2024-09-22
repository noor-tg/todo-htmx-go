package main

import (
	"alnoor/todo-go-htmx"
	server_echo "alnoor/todo-go-htmx/server/echo"
	"crypto/tls"
	"log"
	"net/http"
)

const (
	port = ":3000"
)

func main() {
	serve := server_echo.NewTasksServer(todo.ProductionCfg)
	// Read the embedded certificate and key
	certData, err := todo.Certs.ReadFile("certs/cert.pem")
	if err != nil {
		log.Fatalf("Failed to read embedded cert file: %v", err)
	}

	keyData, err := todo.Certs.ReadFile("certs/key.pem")
	if err != nil {
		log.Fatalf("Failed to read embedded key file: %v", err)
	}

	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}

	// Set up TLS config with the certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	srvr := &http.Server{
		Addr:      port,
		Handler:   serve.Router,
		TLSConfig: tlsConfig,
	}

	log.Println("Starting HTTPS server... https://todo.local:3000")
	err = srvr.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
