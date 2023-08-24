package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/lib"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/scheme"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
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

			resBody := &web.WebResponse[string]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			assert.Equal(t, 400, resBody.Code)
			assert.Equal(t, "failed", resBody.Status)
			assert.Equal(t, "request body is invalid", resBody.Data)

			resetTodolistsDB()
		})
	}
}

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

func TestDeleteTodolistSuccess(t *testing.T) {
	initialDataOne := writeTodolistDB(&domain.Todolist{
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo 1",
	})

	time.Sleep(1 * time.Millisecond)

	initialDataTwo := writeTodolistDB(&domain.Todolist{
		Tags:            []string{"Bar"},
		TodolistMessage: "Initial Todo 2",
	})

	time.Sleep(1 * time.Millisecond)

	initialDataThree := writeTodolistDB(&domain.Todolist{
		Tags:            []string{"Doe"},
		TodolistMessage: "Initial Todo 3",
	})

	time.Sleep(1 * time.Millisecond)

	router := setupRouterTest()

	target := fmt.Sprintf(todolistByIdPath, initialDataOne.Id)
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
	assert.ElementsMatch(t, []domain.Todolist{*initialDataTwo, *initialDataThree}, todolistDB.Todolists)

	resetTodolistsDB()
}

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
