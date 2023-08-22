//go:build wireinject
// +build wireinject

package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/app"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"
)

var todolistRepoSet = wire.NewSet(repository.NewTodolistRepositoryImpl, wire.Bind(new(repository.TodolistRepository), new(*repository.TodolistRepositoryImpl)))

var todolistServiceSet = wire.NewSet(service.NewTodolistServiceImpl, wire.Bind(new(service.TodolistService), new(*service.TodolistServiceImpl)))

var todolistControllerSet = wire.NewSet(controller.NewTodolistControllerImpl, wire.Bind(new(controller.TodolistController), new(*controller.TodolistControllerImpl)))

func InitializeServerTest(dbPath string) *httprouter.Router {
	wire.Build(validator.New, todolistRepoSet, todolistServiceSet, todolistControllerSet, app.NewRouter)
	return nil
}
