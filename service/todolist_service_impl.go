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

func NewTodolistServiceImpl(todolistRepository repository.TodolistRepository, dbPath string, validate *validator.Validate) TodolistService {
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

	foundedTodolist, err := service.TodolistRepository.FindById(service.DBPath, request.Id)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	foundedTodolist.Tags = request.Tags
	foundedTodolist.Done = request.Done
	foundedTodolist.TodolistMessage = request.TodolistMessage

	foundedTodolist = service.TodolistRepository.Update(service.DBPath, foundedTodolist)

	return helper.ToTodolistResponse(foundedTodolist)
}

func (service *TodolistServiceImpl) Delete(todolistIdParam string) {
	foundedTodolist, err := service.TodolistRepository.FindById(service.DBPath, todolistIdParam)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.TodolistRepository.Delete(service.DBPath, foundedTodolist)
}

func (service *TodolistServiceImpl) FindAll() []web.TodolistResponse {
	todolists := service.TodolistRepository.FindAll(service.DBPath)
	return helper.ToTodolistsResponse(todolists)
}
