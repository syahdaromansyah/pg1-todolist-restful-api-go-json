package controller

import (
	"net/http"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"

	"github.com/julienschmidt/httprouter"
)

type TodolistControllerImpl struct {
	TodolistService service.TodolistService
}

func NewTodolistControllerImpl(todolistService service.TodolistService) *TodolistControllerImpl {
	return &TodolistControllerImpl{
		TodolistService: todolistService,
	}
}

func (controller *TodolistControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	todolistCreateRequest := web.TodolistCreateRequest{}
	helper.ReadFromRequestBody(request, &todolistCreateRequest)

	todolistResponse := controller.TodolistService.Create(todolistCreateRequest)

	webResponse := web.WebResponse[web.TodolistResponse]{
		Code:   http.StatusCreated,
		Status: "success",
		Data:   todolistResponse,
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusCreated)
}

func (controller *TodolistControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	todolistUpdateRequest := web.TodolistUpdateRequest{}
	helper.ReadFromRequestBody(request, &todolistUpdateRequest)

	todolistId := params.ByName("todolistId")
	todolistUpdateRequest.Id = todolistId

	todolistResponse := controller.TodolistService.Update(todolistUpdateRequest)
	webResponse := web.WebResponse[web.TodolistResponse]{
		Code:   http.StatusOK,
		Status: "success",
		Data:   todolistResponse,
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusOK)
}

func (controller *TodolistControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	todolistId := params.ByName("todolistId")
	controller.TodolistService.Delete(todolistId)
	webResponse := web.WebResponse[struct{}]{
		Code:   http.StatusOK,
		Status: "success",
		Data:   struct{}{},
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusOK)
}

func (controller *TodolistControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	todolistsResponse := controller.TodolistService.FindAll()
	webResponse := web.WebResponse[[]web.TodolistResponse]{
		Code:   http.StatusOK,
		Status: "success",
		Data:   todolistsResponse,
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusOK)
}
