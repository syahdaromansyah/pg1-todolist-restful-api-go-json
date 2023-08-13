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

func InitializeServerTest(dbPath string) *httprouter.Router {
	wire.Build(validator.New, repository.NewTodolistRepositoryImpl, service.NewTodolistServiceImpl, controller.NewTodolistControllerImpl, app.NewRouter)
	return nil
}
