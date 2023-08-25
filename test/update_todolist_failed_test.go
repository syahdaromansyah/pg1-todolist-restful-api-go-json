package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func TestUpdateTodolistFailed(t *testing.T) {
	initialData := writeTodolistDB(&domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(1 * time.Millisecond)

	tableTests := []domain.Todolist{
		{
			Id:              "notfound",
			Done:            false,
			Tags:            []string{"Not Found"},
			TodolistMessage: "Not Found",
		},
		{
			Done:            false,
			Tags:            []string{""},
			TodolistMessage: "",
		},
		{
			Done:            false,
			Tags:            []string{"", ""},
			TodolistMessage: "Updated Todo 2",
		},
		{
			Done:            true,
			Tags:            []string{"", "Bax", ""},
			TodolistMessage: "Updated Todo 3",
		},
		{
			Done:            true,
			Tags:            []string{"", "Bax", ""},
			TodolistMessage: "",
		},
		{
			Done:            true,
			Tags:            []string{"Bar", "Baz"},
			TodolistMessage: "",
		},
	}

	for _, todolistUpdateReq := range tableTests {
		t.Run(todolistUpdateReq.TodolistMessage, func(t *testing.T) {
			router := setupRouterTest()
			payload := &web.TodolistUpdateRequest{
				Done:            todolistUpdateReq.Done,
				Tags:            todolistUpdateReq.Tags,
				TodolistMessage: todolistUpdateReq.TodolistMessage,
			}

			payloadBytes, err := json.Marshal(payload)
			helper.DoPanicIfError(err)

			reqBody := strings.NewReader(string(payloadBytes))
			var target string

			if todolistUpdateReq.Id == "notfound" {
				target = fmt.Sprintf(todolistByIdPath, todolistUpdateReq.Id)
			} else {
				target = fmt.Sprintf(todolistByIdPath, initialData.Id)
			}

			httpReq := httptest.NewRequest(http.MethodPut, target, reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()

			if todolistUpdateReq.Id == "notfound" {
				assert.Equal(t, 404, response.StatusCode)
			} else {
				assert.Equal(t, 400, response.StatusCode)
			}

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &web.WebResponse[string]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			if todolistUpdateReq.Id == "notfound" {
				assert.Equal(t, 404, resBody.Code)
				assert.Equal(t, "todolist is not found", resBody.Data)
			} else {
				assert.Equal(t, 400, resBody.Code)
				assert.Equal(t, "request body is invalid", resBody.Data)
			}

			assert.Equal(t, "failed", resBody.Status)
		})
	}

	resetTodolistsDB()
}
