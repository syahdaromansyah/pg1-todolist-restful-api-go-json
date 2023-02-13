package main

import (
	"net/http"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/app"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"

	"github.com/go-playground/validator/v10"
)

func main() {
	dbPath := "./databases/todolist.json"
	validate := validator.New()
	todolistRepository := repository.NewTodolistRepositoryImpl()
	todolistService := service.NewTodolistServiceImpl(todolistRepository, dbPath, validate)
	todolistController := controller.NewTodolistControllerImpl(todolistService)
	httpRouter := app.NewRouter(todolistController)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: httpRouter,
	}

	err := server.ListenAndServe()
	helper.DoPanicIfError(err)
}
