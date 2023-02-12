package helper

import (
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func ToTodolistResponse(todolist domain.Todolist) web.TodolistResponse {
	return web.TodolistResponse{
		Id:              todolist.Id,
		Done:            todolist.Done,
		Tags:            todolist.Tags,
		TodolistMessage: todolist.TodolistMessage,
		CreatedAt:       todolist.CreatedAt,
		UpdatedAt:       todolist.UpdatedAt,
	}
}

func ToTodolistsResponse(todolists []domain.Todolist) []web.TodolistResponse {
	todolistsResponse := []web.TodolistResponse{}

	for _, todolist := range todolists {
		todolistsResponse = append(todolistsResponse, ToTodolistResponse(todolist))
	}

	return todolistsResponse
}
