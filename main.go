package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
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

	err := server.ListenAndServe()
	helper.DoPanicIfError(err)
}
