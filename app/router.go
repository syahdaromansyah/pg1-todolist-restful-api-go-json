package app

import (
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(todolistController controller.TodolistController) *httprouter.Router {
	httpRouter := httprouter.New()

	httpRouter.GET("/api/todolists", todolistController.FindAll)
	httpRouter.POST("/api/todolists", todolistController.Create)
	httpRouter.PUT("/api/todolists/:todolistId", todolistController.Update)
	httpRouter.DELETE("/api/todolists/:todolistId", todolistController.Delete)

	httpRouter.PanicHandler = helper.HttpRouterPanicHandler

	return httpRouter
}
