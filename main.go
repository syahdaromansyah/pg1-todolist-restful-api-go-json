package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
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
			helper.WriteLogToFile(func() {
				helper.Logger.Errorf("HTTP server shutdown: %v", err)
			})
		} else {
			helper.WriteLogToFile(func() {
				helper.Logger.Info("HTTP server shutdown gracefully")
			})
		}

		close(idleConnsClosed)
	}()

	helper.WriteLogToFile(func() {
		helper.Logger.Infof("Listening HTTP server on %s", server.Addr)
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		helper.WriteLogToFile(func() {
			helper.Logger.Fatalf("HTTP server ListenAndServe: %v", err)
		})
	}

	<-idleConnsClosed
}
