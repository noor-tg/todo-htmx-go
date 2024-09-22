package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/noor-tg/todo-htmx-go"
	server_echo "github.com/noor-tg/todo-htmx-go/server/echo"
)

const (
	port = ":3000"
)

func https() *tls.Config {
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

	return tlsConfig

}

func main() {
	serve := server_echo.NewTasksServer(todo.ProductionCfg)

	tlsConfig := https()

	srvr := &http.Server{
		Addr:      port,
		Handler:   serve.Router,
		TLSConfig: tlsConfig,
	}

	log.Println("Starting HTTPS server... https://todo.local:3000")

	if err := srvr.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("ListenAndServe: %v\n", err)
	}
}
