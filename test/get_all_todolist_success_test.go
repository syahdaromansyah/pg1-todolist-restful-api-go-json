package test

import (
	"encoding/json"
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

func TestGetAllTodolistSuccess(t *testing.T) {
	tableTests := [][]domain.Todolist{
		{
			{
				Tags:            []string{"Foo"},
				TodolistMessage: "Initial Todo A",
			},
			{
				Tags:            []string{"Bar"},
				TodolistMessage: "Initial Todo B",
			},
		},
		{
			{
				Tags:            []string{"Foo", "Bar"},
				TodolistMessage: "Initial Todo AA",
			},
			{
				Tags:            []string{"Bar", "Ray"},
				TodolistMessage: "Initial Todo AB",
			},
			{
				Tags:            []string{"Far"},
				TodolistMessage: "Initial Todo AC",
			},
		},
	}

	for _, todolistReqs := range tableTests {
		for _, todolistReq := range todolistReqs {
			writeTodolistDB(&todolistReq)
			time.Sleep(1 * time.Millisecond)
		}

		router := setupRouterTest()

		httpReq := httptest.NewRequest(http.MethodGet, todolistsPath, nil)
		httpReq.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httpReq)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		resBodyBytes, err := io.ReadAll(response.Body)
		helper.DoPanicIfError(err)

		resBody := &web.WebResponse[[]domain.Todolist]{}

		err = json.Unmarshal(resBodyBytes, resBody)
		helper.DoPanicIfError(err)

		todolistDB := readTodolistDB()

		assert.Equal(t, 200, resBody.Code)
		assert.Equal(t, "success", resBody.Status)
		assert.Equal(t, len(todolistReqs), len(resBody.Data))
		assert.ElementsMatch(t, todolistDB.Todolists, resBody.Data)

		resetTodolistsDB()
	}
}
