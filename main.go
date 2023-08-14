package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewServer(httpRouter *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    "localhost:8080",
		Handler: httpRouter,
	}
}

func main() {
	dbPath := "./databases/todolist.json"
	server := InitializeServer(dbPath)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: HTTP server ListenAndServe: %v\n", err)
	}
}
