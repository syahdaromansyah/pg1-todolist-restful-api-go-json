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

func TestDeleteTodolistFailed(t *testing.T) {
	initialData := writeTodolistDB(&domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(1 * time.Millisecond)

	router := setupRouterTest()

	target := fmt.Sprintf(todolistByIdPath, "notfound")
	httpReq := httptest.NewRequest(http.MethodDelete, target, nil)
	httpReq.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httpReq)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	resBodyBytes, err := io.ReadAll(response.Body)
	helper.DoPanicIfError(err)

	resBody := &web.WebResponse[string]{}

	err = json.Unmarshal(resBodyBytes, resBody)
	helper.DoPanicIfError(err)

	assert.Equal(t, 404, resBody.Code)
	assert.Equal(t, "failed", resBody.Status)
	assert.Equal(t, "todolist is not found", resBody.Data)

	todolistDB := readTodolistDB()

	assert.Equal(t, uint(1), todolistDB.Total)
	assert.Equal(t, 1, len(todolistDB.Todolists))
	assert.NotEqual(t, 0, len(todolistDB.Todolists))
	assert.ElementsMatch(t, []domain.Todolist{*initialData}, todolistDB.Todolists)

	resetTodolistsDB()
}
