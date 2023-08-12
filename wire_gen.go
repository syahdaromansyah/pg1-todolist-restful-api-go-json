// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/app"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"
	"net/http"
)

// Injectors from injector.go:

func InitializeServer(dbPath string) *http.Server {
	todolistRepository := repository.NewTodolistRepositoryImpl()
	validate := validator.New()
	todolistService := service.NewTodolistServiceImpl(todolistRepository, dbPath, validate)
	todolistController := controller.NewTodolistControllerImpl(todolistService)
	router := app.NewRouter(todolistController)
	server := app.NewServer(router)
	return server
}
