package service

import (
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/exception"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"

	"github.com/go-playground/validator/v10"
)

type TodolistServiceImpl struct {
	TodolistRepository repository.TodolistRepository
	DBPath             string
	Validate           *validator.Validate
}

func NewTodolistServiceImpl(todolistRepository repository.TodolistRepository, dbPath string, validate *validator.Validate) *TodolistServiceImpl {
	return &TodolistServiceImpl{
		TodolistRepository: todolistRepository,
		DBPath:             dbPath,
		Validate:           validate,
	}
}

func (service *TodolistServiceImpl) Create(request web.TodolistCreateRequest) web.TodolistResponse {
	err := service.Validate.Struct(request)
	helper.DoPanicIfError(err)

	todolistRequest := domain.Todolist{
		Tags:            request.Tags,
		TodolistMessage: request.TodolistMessage,
	}

	todolistSavedData := service.TodolistRepository.Save(service.DBPath, todolistRequest)

	return helper.ToTodolistResponse(todolistSavedData)
}

func (service *TodolistServiceImpl) Update(request web.TodolistUpdateRequest) web.TodolistResponse {
	err := service.Validate.Struct(request)
	helper.DoPanicIfError(err)

	updatedTodolist, err := service.TodolistRepository.Update(service.DBPath, domain.Todolist{
		Id:              request.Id,
		Done:            request.Done,
		Tags:            request.Tags,
		TodolistMessage: request.TodolistMessage,
	})

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToTodolistResponse(updatedTodolist)
}

func (service *TodolistServiceImpl) Delete(todolistIdParam string) {
	if err := service.TodolistRepository.Delete(service.DBPath, todolistIdParam); err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *TodolistServiceImpl) FindAll() []web.TodolistResponse {
	todolists := service.TodolistRepository.FindAll(service.DBPath)
	return helper.ToTodolistsResponse(todolists)
}
