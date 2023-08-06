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

	if err := json.Unmarshal(todolistsJsonBytes, todolistsDB); err != nil {
		helper.DoPanicIfError(err)
	}

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

	if err := os.WriteFile(dbPath, marshalledTodolistDB, 0644); err != nil {
		helper.DoPanicIfError(err)
	}

	todolistRequest.Id = newId
	todolistRequest.CreatedAt = createdAt
	todolistRequest.UpdatedAt = *updatedAt

	return todolistRequest
}

func (repository *TodolistRepositoryImpl) Update(dbPath string, todolistRequest domain.Todolist) domain.Todolist {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsJsonMap := map[string]any{}

	jsonUnmarshallErr := json.Unmarshal(todolistsJsonBytes, &todolistsJsonMap)
	helper.DoPanicIfError(jsonUnmarshallErr)

	todolists := todolistsJsonMap["todolists"].([]any)

	for idx, todolist := range todolists {
		todolistId := todolist.(map[string]any)["id"].(string)
		if todolistRequest.Id == todolistId {
			updatedAt := time.Now().Format(time.RFC3339)

			todolistsJsonMap["todolists"].([]any)[idx].(map[string]any)["tags"] = todolistRequest.Tags
			todolistsJsonMap["todolists"].([]any)[idx].(map[string]any)["done"] = todolistRequest.Done
			todolistsJsonMap["todolists"].([]any)[idx].(map[string]any)["todolistMessage"] = todolistRequest.TodolistMessage
			todolistsJsonMap["todolists"].([]any)[idx].(map[string]any)["updatedAt"] = updatedAt

			todolistRequest.UpdatedAt = updatedAt
			break
		}
	}

	jsonMarshalledBytes, err := json.Marshal(todolistsJsonMap)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, jsonMarshalledBytes, 0644)
	helper.DoPanicIfError(writeFileErr)

	return todolistRequest
}

func (repository *TodolistRepositoryImpl) Delete(dbPath string, todolistRequest domain.Todolist) {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsJsonMap := map[string]any{}

	jsonUnmarshallErr := json.Unmarshal(todolistsJsonBytes, &todolistsJsonMap)
	helper.DoPanicIfError(jsonUnmarshallErr)

	todolists := todolistsJsonMap["todolists"].([]any)
	todolistsJsonMap["todolists"] = []any{}

	for _, todolist := range todolists {
		todolistId := todolist.(map[string]any)["id"].(string)
		if todolistRequest.Id == todolistId {
			todolistsJsonMap["total"] = int(todolistsJsonMap["total"].(float64)) - 1
			continue
		}

		todolistsJsonMap["todolists"] = append(todolistsJsonMap["todolists"].([]any), todolist)
	}

	jsonMarshalledBytes, err := json.Marshal(todolistsJsonMap)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, jsonMarshalledBytes, 0644)
	helper.DoPanicIfError(writeFileErr)
}

func (repository *TodolistRepositoryImpl) FindById(dbPath, todolistIdParam string) (domain.Todolist, error) {
	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsJsonMap := map[string]any{}

	jsonUnmarshallErr := json.Unmarshal(todolistsJsonBytes, &todolistsJsonMap)
	helper.DoPanicIfError(jsonUnmarshallErr)

	todolists := todolistsJsonMap["todolists"].([]any)
	foundedTodolist := domain.Todolist{}

	for _, todolist := range todolists {
		todolistId := todolist.(map[string]any)["id"].(string)

		if todolistIdParam == todolistId {
			foundedTodolist.Id = todolistId
			foundedTodolist.Done = todolist.(map[string]any)["done"].(bool)
			foundedTodolist.Tags = []string{}
			foundedTodolist.TodolistMessage = todolist.(map[string]any)["todolistMessage"].(string)
			foundedTodolist.CreatedAt = todolist.(map[string]any)["createdAt"].(string)
			foundedTodolist.UpdatedAt = todolist.(map[string]any)["updatedAt"].(string)

			for _, tag := range todolist.(map[string]any)["tags"].([]any) {
				foundedTodolist.Tags = append(foundedTodolist.Tags, tag.(string))
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
