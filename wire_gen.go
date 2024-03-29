// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/app"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"
	"net/http"
)

// Injectors from injector.go:

func InitializeServer(dbPath string) *http.Server {
	todolistRepositoryImpl := repository.NewTodolistRepositoryImpl()
	validate := validator.New()
	todolistServiceImpl := service.NewTodolistServiceImpl(todolistRepositoryImpl, dbPath, validate)
	todolistControllerImpl := controller.NewTodolistControllerImpl(todolistServiceImpl)
	router := app.NewRouter(todolistControllerImpl)
	server := app.NewServer(router)
	return server
}

// injector.go:

var todolistRepoSet = wire.NewSet(repository.NewTodolistRepositoryImpl, wire.Bind(new(repository.TodolistRepository), new(*repository.TodolistRepositoryImpl)))

var todolistServiceSet = wire.NewSet(service.NewTodolistServiceImpl, wire.Bind(new(service.TodolistService), new(*service.TodolistServiceImpl)))

var todolistControllerSet = wire.NewSet(controller.NewTodolistControllerImpl, wire.Bind(new(controller.TodolistController), new(*controller.TodolistControllerImpl)))
