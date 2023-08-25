package test

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/lib"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/scheme"
)

const dbPath = "../databases/todolist.json"
const todolistsPath = "http://localhost:8080/api/todolists"
const todolistByIdPath = "http://localhost:8080/api/todolists/%s"

func setupRouterTest() http.Handler {
	httpRouter := InitializeServerTest(dbPath)
	return httpRouter
}

func resetTodolistsDB() {
	err := os.WriteFile(dbPath, []byte(`{ "todolists": [], "total": 0 }`+"\n"), 0644)
	helper.DoPanicIfError(err)
}

func writeTodolistDB(todolistReq *domain.Todolist) *domain.Todolist {
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
		Tags:            todolistReq.Tags,
		TodolistMessage: todolistReq.TodolistMessage,
		CreatedAt:       createdAt,
		UpdatedAt:       *updatedAt,
	})

	marshalledTodolistDB, err := json.Marshal(todolistsDB)
	helper.DoPanicIfError(err)

	writeFileErr := os.WriteFile(dbPath, marshalledTodolistDB, 0644)
	helper.DoPanicIfError(writeFileErr)

	todolistReq.Id = newId
	todolistReq.CreatedAt = createdAt
	todolistReq.UpdatedAt = *updatedAt

	return todolistReq
}

func readTodolistDB() *scheme.TodolistDB {
	todolistJsonBytes, err := os.ReadFile(dbPath)
	helper.DoPanicIfError(err)

	todolistDB := &scheme.TodolistDB{}
	err = json.Unmarshal(todolistJsonBytes, todolistDB)
	helper.DoPanicIfError(err)

	return todolistDB
}
