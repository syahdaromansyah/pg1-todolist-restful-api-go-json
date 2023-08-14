package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	dbPath := "./databases/todolist.json"
	server := InitializeServer(dbPath)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: HTTP server ListenAndServe: %v\n", err)
	}
}
