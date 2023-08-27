package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func TestCreateTodolistFailed(t *testing.T) {
	tableTests := []struct {
		Data web.TodolistCreateRequest
	}{
		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{},
				TodolistMessage: "",
			},
		},
		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{""},
				TodolistMessage: "This is a Todo",
			},
		},
		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{"Programming", "", "Technology"},
				TodolistMessage: "Learn Go Unit Test",
			},
		},

		{
			Data: web.TodolistCreateRequest{Tags: []string{"Cloud", "Virtual Machine"},
				TodolistMessage: ""},
		},

		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{"", "", ""},
				TodolistMessage: "",
			},
		},
	}

	for _, test := range tableTests {
		t.Run(test.Data.TodolistMessage, func(t *testing.T) {
			router := setupRouterTest()
			payload := &web.TodolistCreateRequest{
				Tags:            test.Data.Tags,
				TodolistMessage: test.Data.TodolistMessage,
			}

			payloadBytes, err := json.Marshal(payload)
			helper.DoPanicIfError(err)

			reqBody := strings.NewReader(string(payloadBytes))
			httpReq := httptest.NewRequest(http.MethodPost, todolistsPath, reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()
			assert.Equal(t, 400, response.StatusCode)

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &web.WebResponse[struct{}]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			assert.Equal(t, 400, resBody.Code)
			assert.Equal(t, "failed", resBody.Status)
			assert.Equal(t, struct{}{}, resBody.Data)

			// Ensure invalid todolist request is not added to DB
			todolistDB := readTodolistDB()

			assert.Equal(t, uint(0), todolistDB.Total)
			assert.Equal(t, 0, len(todolistDB.Todolists))

			resetTodolistsDB()
		})
	}
}
