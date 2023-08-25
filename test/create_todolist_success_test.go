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
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func TestCreateTodolistSuccess(t *testing.T) {
	type todolistReqTest struct {
		TodolistReq       web.TodolistCreateRequest
		ExpectedTotalTags int
	}

	tableTests := []struct {
		TodolistReqsTest      []todolistReqTest
		ExpectedTotalTodolist uint
	}{
		{
			TodolistReqsTest: []todolistReqTest{
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{},
						TodolistMessage: "This is a Todo",
					},
					ExpectedTotalTags: 0,
				},
			},
			ExpectedTotalTodolist: 1,
		},
		{
			TodolistReqsTest: []todolistReqTest{
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Programming"},
						TodolistMessage: "Learn Go Unit Test",
					},
					ExpectedTotalTags: 1,
				},
			},
			ExpectedTotalTodolist: 1,
		},
		{
			TodolistReqsTest: []todolistReqTest{
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Cloud", "Virtual Machine"},
						TodolistMessage: "Learn Virtual Machine in Cloud",
					},
					ExpectedTotalTags: 2,
				},
			},
			ExpectedTotalTodolist: 1,
		},
		{
			TodolistReqsTest: []todolistReqTest{
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Programming", "OOP", "Async"},
						TodolistMessage: "Learn Java Async",
					},
					ExpectedTotalTags: 3,
				},
			},
			ExpectedTotalTodolist: 1,
		},
		{
			TodolistReqsTest: []todolistReqTest{
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{},
						TodolistMessage: "This is a Todo",
					},
					ExpectedTotalTags: 0,
				},
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Programming"},
						TodolistMessage: "Learn Go Unit Test",
					},
					ExpectedTotalTags: 1,
				},
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Cloud", "Virtual Machine"},
						TodolistMessage: "Learn Virtual Machine in Cloud",
					},
					ExpectedTotalTags: 2,
				},
				{
					TodolistReq: web.TodolistCreateRequest{
						Tags:            []string{"Programming", "OOP", "Async"},
						TodolistMessage: "Learn Java Async",
					},
					ExpectedTotalTags: 3,
				},
			},
			ExpectedTotalTodolist: 4,
		},
	}

	for _, test := range tableTests {
		todolistsRes := []domain.Todolist{}

		for _, todolistReqTest := range test.TodolistReqsTest {
			todolistReqTags := todolistReqTest.TodolistReq.Tags
			todolistReqMessage := todolistReqTest.TodolistReq.TodolistMessage

			t.Run(todolistReqMessage, func(t *testing.T) {
				router := setupRouterTest()
				payload := &web.TodolistCreateRequest{
					Tags:            todolistReqTags,
					TodolistMessage: todolistReqMessage,
				}

				payloadBytes, err := json.Marshal(payload)
				helper.DoPanicIfError(err)

				reqBody := strings.NewReader(string(payloadBytes))
				httpReq := httptest.NewRequest(http.MethodPost, todolistsPath, reqBody)
				httpReq.Header.Add("Content-Type", "application/json")

				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, httpReq)

				response := recorder.Result()
				assert.Equal(t, 201, response.StatusCode)

				resBodyBytes, err := io.ReadAll(response.Body)
				helper.DoPanicIfError(err)

				resBody := &web.WebResponse[domain.Todolist]{}

				err = json.Unmarshal(resBodyBytes, resBody)
				helper.DoPanicIfError(err)

				assert.Equal(t, 201, resBody.Code)
				assert.Equal(t, "success", resBody.Status)
				assert.Equal(t, todolistReqTest.ExpectedTotalTags, len(resBody.Data.Tags))
				assert.ElementsMatch(t, todolistReqTags, resBody.Data.Tags)
				assert.Equal(t, todolistReqMessage, resBody.Data.TodolistMessage)
				assert.Equal(t, resBody.Data.UpdatedAt, resBody.Data.CreatedAt)
				assert.NotEmpty(t, resBody.Data.Id)
				assert.False(t, resBody.Data.Done)

				todolistsRes = append(todolistsRes, domain.Todolist{
					Id:              resBody.Data.Id,
					Done:            resBody.Data.Done,
					Tags:            resBody.Data.Tags,
					TodolistMessage: resBody.Data.TodolistMessage,
					CreatedAt:       resBody.Data.CreatedAt,
					UpdatedAt:       resBody.Data.UpdatedAt,
				})
			})
		}

		todolistDB := readTodolistDB()

		assert.ElementsMatch(t, todolistDB.Todolists, todolistsRes)
		assert.Equal(t, test.ExpectedTotalTodolist, todolistDB.Total)

		resetTodolistsDB()
	}
}
