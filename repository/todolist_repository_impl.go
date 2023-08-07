package repository

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/lib"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/scheme"
)

type TodolistRepositoryImpl struct{}

func NewTodolistRepositoryImpl() TodolistRepository {
	return &TodolistRepositoryImpl{}
}

func (repository *TodolistRepositoryImpl) Save(dbPath string, todolistRequest domain.Todolist) domain.Todolist {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	newId := lib.GetRandomStdId32()
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := &createdAt

	todolistsDB.Total = todolistsDB.Total + 1
	todolistsDB.Todolists = append(todolistsDB.Todolists, domain.Todolist{
		Id:              newId,
		Done:            false,
		Tags:            todolistRequest.Tags,
		TodolistMessage: todolistRequest.TodolistMessage,
		CreatedAt:       createdAt,
		UpdatedAt:       *updatedAt,
	})

	marshalledTodolistDB, err := json.Marshal(todolistsDB)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, marshalledTodolistDB, 0644)
	helper.DoPanicIfError(writeFileErr)

	todolistRequest.Id = newId
	todolistRequest.CreatedAt = createdAt
	todolistRequest.UpdatedAt = *updatedAt

	return todolistRequest
}

func (repository *TodolistRepositoryImpl) Update(dbPath string, todolistRequest domain.Todolist) domain.Todolist {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	todolists := todolistsDB.Todolists

	for idx, todolist := range todolists {
		if todolistRequest.Id == todolist.Id {
			updatedAt := time.Now().Format(time.RFC3339)
			todolistsDB.Todolists[idx].Tags = todolistRequest.Tags
			todolistsDB.Todolists[idx].Done = todolistRequest.Done
			todolistsDB.Todolists[idx].TodolistMessage = todolistRequest.TodolistMessage
			todolistsDB.Todolists[idx].UpdatedAt = todolistRequest.UpdatedAt

			todolistRequest.UpdatedAt = updatedAt
			break
		}
	}

	marshalledTodolistDB, err := json.Marshal(todolistsDB)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, marshalledTodolistDB, 0644)
	helper.DoPanicIfError(writeFileErr)

	return todolistRequest
}

func (repository *TodolistRepositoryImpl) Delete(dbPath string, todolistRequest domain.Todolist) {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	todolists := todolistsDB.Todolists
	todolistsDB.Todolists = []domain.Todolist{}

	for _, todolist := range todolists {
		if todolistRequest.Id == todolist.Id {
			todolistsDB.Total = todolistsDB.Total - 1
			continue
		}

		todolistsDB.Todolists = append(todolistsDB.Todolists, todolist)
	}

	jsonMarshalledBytes, err := json.Marshal(todolistsDB)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, jsonMarshalledBytes, 0644)
	helper.DoPanicIfError(writeFileErr)
}

func (repository *TodolistRepositoryImpl) FindById(dbPath, todolistIdParam string) (domain.Todolist, error) {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	todolists := todolistsDB.Todolists
	foundedTodolist := domain.Todolist{}

	for _, todolist := range todolists {
		if todolistIdParam == todolist.Id {
			foundedTodolist.Id = todolist.Id
			foundedTodolist.Done = todolist.Done
			foundedTodolist.TodolistMessage = todolist.TodolistMessage
			foundedTodolist.CreatedAt = todolist.CreatedAt
			foundedTodolist.UpdatedAt = todolist.UpdatedAt

			if len(todolist.Tags) == 0 {
				foundedTodolist.Tags = []string{}
			} else {
				foundedTodolist.Tags = append(foundedTodolist.Tags, todolist.Tags...)
			}

			return foundedTodolist, nil
		}
	}

	return foundedTodolist, errors.New("todolist is not found")
}

func (repository *TodolistRepositoryImpl) FindAll(dbPath string) []domain.Todolist {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsJsonMap := map[string]any{}

	jsonUnmarshallErr := json.Unmarshal(todolistsJsonBytes, &todolistsJsonMap)
	helper.DoPanicIfError(jsonUnmarshallErr)

	todolistsJson := todolistsJsonMap["todolists"].([]any)
	todolists := []domain.Todolist{}

	for _, todolist := range todolistsJson {
		todolistData := domain.Todolist{
			Id:              todolist.(map[string]any)["id"].(string),
			Done:            todolist.(map[string]any)["done"].(bool),
			Tags:            []string{},
			TodolistMessage: todolist.(map[string]any)["todolistMessage"].(string),
			CreatedAt:       todolist.(map[string]any)["createdAt"].(string),
			UpdatedAt:       todolist.(map[string]any)["updatedAt"].(string),
		}

		for _, tag := range todolist.(map[string]any)["tags"].([]any) {
			todolistData.Tags = append(todolistData.Tags, tag.(string))
		}

		todolists = append(todolists, todolistData)
	}

	return todolists
}
