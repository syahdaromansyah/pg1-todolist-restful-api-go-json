package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/lib"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/scheme"
)

var mutex = &sync.Mutex{}

type TodolistRepositoryImpl struct{}

func NewTodolistRepositoryImpl() *TodolistRepositoryImpl {
	return &TodolistRepositoryImpl{}
}

func (repository *TodolistRepositoryImpl) Save(dbPath string, todolistRequest domain.Todolist) domain.Todolist {
	mutex.Lock()
	defer mutex.Unlock()

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

func (repository *TodolistRepositoryImpl) Update(dbPath string, todolistRequest domain.Todolist) (domain.Todolist, error) {
	mutex.Lock()
	defer mutex.Unlock()

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
			todolistsDB.Todolists[idx].UpdatedAt = updatedAt

			todolistRequest.CreatedAt = todolist.CreatedAt
			todolistRequest.UpdatedAt = updatedAt

			marshalledTodolistDB, err := json.Marshal(todolistsDB)
			helper.DoPanicIfError(err)

			writeFileErr := os.WriteFile(dbPath, marshalledTodolistDB, 0644)
			helper.DoPanicIfError(writeFileErr)

			return todolistRequest, nil
		}
	}

	return todolistRequest, errors.New("todolist is not found")
}

func (repository *TodolistRepositoryImpl) Delete(dbPath string, todolistIdParam string) error {
	mutex.Lock()
	defer mutex.Unlock()

	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	todolists := todolistsDB.Todolists
	todolistsDB.Todolists = []domain.Todolist{}
	todolistIsFounded := false

	for _, todolist := range todolists {
		if todolistIdParam == todolist.Id {
			todolistsDB.Total = todolistsDB.Total - 1
			todolistIsFounded = true
			continue
		}

		todolistsDB.Todolists = append(todolistsDB.Todolists, todolist)
	}

	if todolistIsFounded {
		jsonMarshalledBytes, err := json.Marshal(todolistsDB)
		helper.DoPanicIfError(err)

		writeFileErr := os.WriteFile(dbPath, jsonMarshalledBytes, 0644)
		helper.DoPanicIfError(writeFileErr)

		return nil
	}

	return errors.New("todolist is not found")
}

func (repository *TodolistRepositoryImpl) FindAll(dbPath string) []domain.Todolist {
	mutex.Lock()
	defer mutex.Unlock()

	todolistsJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistsDB := &scheme.TodolistDB{}

	unMarshallErr := json.Unmarshal(todolistsJsonBytes, todolistsDB)
	helper.DoPanicIfError(unMarshallErr)

	todolistsJson := todolistsDB.Todolists
	todolists := []domain.Todolist{}

	for _, todolist := range todolistsJson {
		todolistData := domain.Todolist{
			Id:              todolist.Id,
			Done:            todolist.Done,
			TodolistMessage: todolist.TodolistMessage,
			CreatedAt:       todolist.CreatedAt,
			UpdatedAt:       todolist.UpdatedAt,
		}

		if len(todolist.Tags) == 0 {
			todolistData.Tags = []string{}
		} else {
			todolistData.Tags = append(todolistData.Tags, todolist.Tags...)
		}

		todolists = append(todolists, todolistData)
	}

	return todolists
}
