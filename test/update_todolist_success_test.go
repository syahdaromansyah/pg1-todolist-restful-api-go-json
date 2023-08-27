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

func TestUpdateTodolistSuccess(t *testing.T) {
	var initTodolistData []domain.Todolist

	for _, initTodolist := range []domain.Todolist{
		{
			Tags:            []string{"Foo"},
			TodolistMessage: "Initial Todo 1",
		},
		{
			Tags:            []string{"Bar"},
			TodolistMessage: "Initial Todo 2",
		}, {
			Tags:            []string{"Doo", "Goo"},
			TodolistMessage: "Initial Todo 3",
		},
	} {
		initTodolistData = append(initTodolistData, *writeTodolistDB(&initTodolist))

		time.Sleep(1 * time.Millisecond)
	}

	tableTests := []struct {
		Done            bool
		Tags            []string
		TodolistMessage string
	}{
		{
			Done:            false,
			Tags:            []string{},
			TodolistMessage: "Updated Todo 1",
		},
		{
			Done:            false,
			Tags:            []string{"Baz"},
			TodolistMessage: "Updated Todo 2",
		},
		{
			Done:            true,
			Tags:            []string{},
			TodolistMessage: "Updated Todo 3",
		},
		{
			Done:            true,
			Tags:            []string{"Bar", "Baz"},
			TodolistMessage: "Updated Todo 4",
		},
	}

	selectedTodolist := initTodolistData[0]

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
			target := fmt.Sprintf(todolistByIdPath, selectedTodolist.Id)
			httpReq := httptest.NewRequest(http.MethodPut, target, reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()
			assert.Equal(t, 200, response.StatusCode)

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &web.WebResponse[domain.Todolist]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			assert.Equal(t, 200, resBody.Code)
			assert.Equal(t, "success", resBody.Status)

			todolistDB := readTodolistDB()

			for _, todolistInDB := range todolistDB.Todolists {
				// Ensure the selected todolist is updated
				if todolistInDB.Id == selectedTodolist.Id {
					assert.Equal(t, len(todolistUpdateReq.Tags), len(resBody.Data.Tags))
					assert.Equal(t, len(todolistUpdateReq.Tags), len(todolistInDB.Tags))

					assert.Equal(t, todolistUpdateReq.Tags, resBody.Data.Tags)
					assert.Equal(t, todolistUpdateReq.Tags, todolistInDB.Tags)

					assert.Equal(t, todolistUpdateReq.TodolistMessage, resBody.Data.TodolistMessage)
					assert.Equal(t, todolistUpdateReq.TodolistMessage, todolistInDB.TodolistMessage)

					assert.Equal(t, resBody.Data.CreatedAt, todolistInDB.CreatedAt)
					assert.Equal(t, resBody.Data.UpdatedAt, todolistInDB.UpdatedAt)

					assert.Equal(t, selectedTodolist.Id, resBody.Data.Id)

					if todolistUpdateReq.Done == false {
						assert.False(t, resBody.Data.Done)
						assert.False(t, todolistInDB.Done)
					} else {
						assert.True(t, resBody.Data.Done)
						assert.True(t, todolistInDB.Done)
					}

					break
				}
			}

			// Ensure another todolist in DB is not updated
			for _, anotherTodolist := range initTodolistData[1:] {
				for _, todolistInDB := range todolistDB.Todolists {
					if anotherTodolist.Id == todolistInDB.Id {
						assert.Equal(t, anotherTodolist.Id, todolistInDB.Id)
						assert.Equal(t, anotherTodolist.Done, todolistInDB.Done)
						assert.ElementsMatch(t, anotherTodolist.Tags, todolistInDB.Tags)
						assert.Equal(t, anotherTodolist.TodolistMessage, todolistInDB.TodolistMessage)
						assert.Equal(t, anotherTodolist.CreatedAt, todolistInDB.CreatedAt)
						assert.Equal(t, anotherTodolist.UpdatedAt, todolistInDB.UpdatedAt)
					}
				}
			}
		})
	}

	resetTodolistsDB()
}
