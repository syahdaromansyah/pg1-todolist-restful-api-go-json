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

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/app"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/controller"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/service"
)

type WebResponseTest[T any] struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

func setupRouterTest() http.Handler {
	dbPath := "../databases/todolist.json"
	validate := validator.New()
	todolistRepository := repository.NewTodolistRepositoryImpl()
	todolistService := service.NewTodolistServiceImpl(todolistRepository, dbPath, validate)
	todolistController := controller.NewTodolistControllerImpl(todolistService)
	httpRouter := app.NewRouter(todolistController)

	return httpRouter
}

func resetTodolistsDB() {
	err := os.WriteFile("../databases/todolist.json", []byte(`{ "todolists": [], "total": 0 }`+"\n"), 0644)
	helper.DoPanicIfError(err)
}

func TestCreateTodolistSuccess(t *testing.T) {
	tableTests := []struct {
		Data            web.TodolistCreateRequest
		ExpectedTagsLen int
	}{
		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{},
				TodolistMessage: "This is a Todo",
			},
			ExpectedTagsLen: 0,
		},
		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{"Programming"},
				TodolistMessage: "Learn Go Unit Test",
			},
			ExpectedTagsLen: 1,
		},

		{
			Data: web.TodolistCreateRequest{Tags: []string{"Cloud", "Virtual Machine"},
				TodolistMessage: "Learn Virtual Machine in Cloud"},

			ExpectedTagsLen: 2,
		},

		{
			Data: web.TodolistCreateRequest{
				Tags:            []string{"Programming", "OOP", "Async"},
				TodolistMessage: "Learn Java Async",
			},
			ExpectedTagsLen: 3,
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
			httpReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/todolists", reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()
			assert.Equal(t, 201, response.StatusCode)

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &WebResponseTest[domain.Todolist]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			assert.Equal(t, 201, resBody.Code)
			assert.Equal(t, "success", resBody.Status)
			assert.Equal(t, test.ExpectedTagsLen, len(resBody.Data.Tags))
			assert.Equal(t, test.Data.TodolistMessage, resBody.Data.TodolistMessage)
			assert.Equal(t, resBody.Data.UpdatedAt, resBody.Data.CreatedAt)
			assert.NotEqual(t, "", resBody.Data.Id)
			assert.False(t, resBody.Data.Done)

			resetTodolistsDB()
		})
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
			httpReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/todolists", reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()
			assert.Equal(t, 400, response.StatusCode)

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &WebResponseTest[string]{}

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
	dbPath := "../databases/todolist.json"
	initialData := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(300 * time.Millisecond)

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

	for _, test := range tableTests {
		t.Run(test.TodolistMessage, func(t *testing.T) {
			router := setupRouterTest()
			payload := &web.TodolistUpdateRequest{
				Done:            test.Done,
				Tags:            test.Tags,
				TodolistMessage: test.TodolistMessage,
			}

			payloadBytes, err := json.Marshal(payload)
			helper.DoPanicIfError(err)

			reqBody := strings.NewReader(string(payloadBytes))
			target := fmt.Sprintf("http://localhost:8080/api/todolists/%s", initialData.Id)
			httpReq := httptest.NewRequest(http.MethodPut, target, reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()
			assert.Equal(t, 200, response.StatusCode)

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &WebResponseTest[domain.Todolist]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			assert.Equal(t, 200, resBody.Code)
			assert.Equal(t, "success", resBody.Status)

			todolist, err := repository.NewTodolistRepositoryImpl().FindById(dbPath, initialData.Id)

			assert.Nil(t, err)

			assert.Equal(t, len(test.Tags), len(resBody.Data.Tags))
			assert.Equal(t, len(test.Tags), len(todolist.Tags))

			assert.Equal(t, test.Tags, resBody.Data.Tags)
			assert.Equal(t, test.Tags, todolist.Tags)

			assert.Equal(t, test.TodolistMessage, resBody.Data.TodolistMessage)
			assert.Equal(t, test.TodolistMessage, todolist.TodolistMessage)

			assert.Equal(t, resBody.Data.CreatedAt, todolist.CreatedAt)
			assert.Equal(t, resBody.Data.UpdatedAt, todolist.UpdatedAt)

			assert.Equal(t, initialData.Id, resBody.Data.Id)

			if test.Done == false {
				assert.False(t, resBody.Data.Done)
				assert.False(t, todolist.Done)
			} else {
				assert.True(t, resBody.Data.Done)
				assert.True(t, todolist.Done)
			}
		})
	}

	resetTodolistsDB()
}

func TestUpdateTodolistFailed(t *testing.T) {
	dbPath := "../databases/todolist.json"
	initialData := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(300 * time.Millisecond)

	tableTests := []struct {
		Id              string
		Done            bool
		Tags            []string
		TodolistMessage string
	}{
		{
			Id:              "notfound",
			Done:            false,
			Tags:            []string{"Sport"},
			TodolistMessage: "Sport Todo",
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

	for _, test := range tableTests {
		t.Run(test.TodolistMessage, func(t *testing.T) {
			router := setupRouterTest()
			payload := &web.TodolistUpdateRequest{
				Done:            test.Done,
				Tags:            test.Tags,
				TodolistMessage: test.TodolistMessage,
			}

			payloadBytes, err := json.Marshal(payload)
			helper.DoPanicIfError(err)

			reqBody := strings.NewReader(string(payloadBytes))
			var target string

			if test.Id == "notfound" {
				target = fmt.Sprintf("http://localhost:8080/api/todolists/%s", test.Id)
			} else {
				target = fmt.Sprintf("http://localhost:8080/api/todolists/%s", initialData.Id)
			}

			httpReq := httptest.NewRequest(http.MethodPut, target, reqBody)
			httpReq.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httpReq)

			response := recorder.Result()

			if test.Id == "notfound" {
				assert.Equal(t, 404, response.StatusCode)
			} else {
				assert.Equal(t, 400, response.StatusCode)
			}

			resBodyBytes, err := io.ReadAll(response.Body)
			helper.DoPanicIfError(err)

			resBody := &WebResponseTest[string]{}

			err = json.Unmarshal(resBodyBytes, resBody)
			helper.DoPanicIfError(err)

			if test.Id == "notfound" {
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
	dbPath := "../databases/todolist.json"
	initialData := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(300 * time.Millisecond)

	router := setupRouterTest()

	target := fmt.Sprintf("http://localhost:8080/api/todolists/%s", initialData.Id)
	httpReq := httptest.NewRequest(http.MethodDelete, target, nil)
	httpReq.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httpReq)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	resBodyBytes, err := io.ReadAll(response.Body)
	helper.DoPanicIfError(err)

	resBody := &WebResponseTest[struct{}]{}

	err = json.Unmarshal(resBodyBytes, resBody)
	helper.DoPanicIfError(err)

	assert.Equal(t, 200, resBody.Code)
	assert.Equal(t, "success", resBody.Status)
	assert.Equal(t, struct{}{}, resBody.Data)

	_, err = repository.NewTodolistRepositoryImpl().FindById(dbPath, initialData.Id)

	assert.NotNil(t, err)
	assert.Equal(t, "todolist is not found", err.Error())

	resetTodolistsDB()
}

func TestDeleteTodolistFailed(t *testing.T) {
	dbPath := "../databases/todolist.json"
	repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Done:            false,
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo",
	})

	time.Sleep(300 * time.Millisecond)

	router := setupRouterTest()

	target := "http://localhost:8080/api/todolists/notfound"
	httpReq := httptest.NewRequest(http.MethodDelete, target, nil)
	httpReq.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httpReq)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	resBodyBytes, err := io.ReadAll(response.Body)
	helper.DoPanicIfError(err)

	resBody := &WebResponseTest[string]{}

	err = json.Unmarshal(resBodyBytes, resBody)
	helper.DoPanicIfError(err)

	assert.Equal(t, 404, resBody.Code)
	assert.Equal(t, "failed", resBody.Status)
	assert.Equal(t, "todolist is not found", resBody.Data)

	resetTodolistsDB()
}

func TestGetAllTodolistSuccess(t *testing.T) {
	tableTests := [][]struct {
		Tags            []string
		TodolistMessage string
	}{
		{},
		{
			{
				Tags:            []string{"Ray"},
				TodolistMessage: "Initial Todo 1",
			},
		},
		{
			{
				Tags:            []string{"Foo"},
				TodolistMessage: "Initial Todo 2",
			},
			{
				Tags:            []string{"Bar", "Baz"},
				TodolistMessage: "Initial Todo 3",
			},
		},
	}

	for _, test := range tableTests {
		for _, data := range test {
			dbPath := "../databases/todolist.json"
			repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
				Tags:            data.Tags,
				TodolistMessage: data.TodolistMessage,
			})
		}

		router := setupRouterTest()

		target := "http://localhost:8080/api/todolists"
		httpReq := httptest.NewRequest(http.MethodGet, target, nil)
		httpReq.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httpReq)

		response := recorder.Result()
		assert.Equal(t, 200, response.StatusCode)

		resBodyBytes, err := io.ReadAll(response.Body)
		helper.DoPanicIfError(err)

		resBody := &WebResponseTest[[]domain.Todolist]{}

		err = json.Unmarshal(resBodyBytes, resBody)
		helper.DoPanicIfError(err)

		assert.Equal(t, 200, resBody.Code)
		assert.Equal(t, "success", resBody.Status)
		assert.Equal(t, len(test), len(resBody.Data))

		resBodyTest := []struct {
			Tags            []string
			TodolistMessage string
		}{}

		for _, resBodyData := range resBody.Data {
			resBodyTest = append(resBodyTest, struct {
				Tags            []string
				TodolistMessage string
			}{
				Tags:            resBodyData.Tags,
				TodolistMessage: resBodyData.TodolistMessage,
			})
		}

		assert.ElementsMatch(t, test, resBodyTest)

		resetTodolistsDB()
	}
}
