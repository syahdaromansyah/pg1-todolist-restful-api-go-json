package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbPath := "./databases/todolist.json"
	server := InitializeServer(dbPath)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigInt := make(chan os.Signal, 1)

		signal.Notify(sigInt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

		<-sigInt

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error: HTTP server shutdown: %v\n", err)
		} else {
			log.Printf("HTTP server shutdown gracefully\n")
		}

		close(idleConnsClosed)
	}()

	log.Printf("Listening HTTP server on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: HTTP server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed
}
