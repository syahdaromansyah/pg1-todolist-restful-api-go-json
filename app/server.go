package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewServer(httpRouter *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    "localhost:8080",
		Handler: httpRouter,
	}
}
