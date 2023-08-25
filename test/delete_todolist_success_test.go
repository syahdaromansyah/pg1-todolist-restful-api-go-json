package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func TestDeleteTodolistSuccess(t *testing.T) {
	var initTodolistData []domain.Todolist

	for _, initTodolist := range []domain.Todolist{
		{
			Tags:            []string{"Foo"},
			TodolistMessage: "Initial Todo 1",
		},
		{
			Tags:            []string{"Bar"},
			TodolistMessage: "Initial Todo 2",
		},
		{
			Tags:            []string{"Doe"},
			TodolistMessage: "Initial Todo 3",
		},
	} {
		initTodolistData = append(initTodolistData, *writeTodolistDB(&initTodolist))
		time.Sleep(1 * time.Millisecond)
	}

	selectedInitTodolist := initTodolistData[0]
	anotherInitTodolist := initTodolistData[1:]

	router := setupRouterTest()

	target := fmt.Sprintf(todolistByIdPath, selectedInitTodolist.Id)
	httpReq := httptest.NewRequest(http.MethodDelete, target, nil)
	httpReq.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httpReq)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	resBodyBytes, err := io.ReadAll(response.Body)
	helper.DoPanicIfError(err)

	resBody := &web.WebResponse[struct{}]{}

	err = json.Unmarshal(resBodyBytes, resBody)
	helper.DoPanicIfError(err)

	assert.Equal(t, 200, resBody.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, struct{}{}, resBody.Data)

	todolistDB := readTodolistDB()

	assert.Equal(t, uint(2), todolistDB.Total)
	assert.Equal(t, 2, len(todolistDB.Todolists))
	assert.NotEqual(t, 3, len(todolistDB.Todolists))
	assert.ElementsMatch(t, anotherInitTodolist, todolistDB.Todolists)

	resetTodolistsDB()
}
