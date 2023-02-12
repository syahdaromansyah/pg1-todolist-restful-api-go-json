package service

import "github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"

type TodolistService interface {
	Create(request web.TodolistCreateRequest) web.TodolistResponse
	Update(request web.TodolistUpdateRequest) web.TodolistResponse
	Delete(todolistIdParam string)
	FindAll() []web.TodolistResponse
}
