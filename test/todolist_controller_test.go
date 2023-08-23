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
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/scheme"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/repository"
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

				todolistsRes = append(todolistsRes, domain.Todolist{
					Id:              resBody.Data.Id,
					Done:            resBody.Data.Done,
					Tags:            resBody.Data.Tags,
					TodolistMessage: resBody.Data.TodolistMessage,
					CreatedAt:       resBody.Data.CreatedAt,
					UpdatedAt:       resBody.Data.UpdatedAt,
				})

				assert.Equal(t, 201, resBody.Code)
				assert.Equal(t, "success", resBody.Status)
				assert.Equal(t, todolistReqTest.ExpectedTotalTags, len(resBody.Data.Tags))
				assert.Equal(t, todolistReqMessage, resBody.Data.TodolistMessage)
				assert.Equal(t, resBody.Data.UpdatedAt, resBody.Data.CreatedAt)
				assert.NotEmpty(t, resBody.Data.Id)
				assert.False(t, resBody.Data.Done)
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
	initialDataOne := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo 1",
	})

	time.Sleep(300 * time.Millisecond)

	initialDataTwo := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Boo"},
		TodolistMessage: "Initial Todo 2",
	})

	time.Sleep(300 * time.Millisecond)

	initialDataThree := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Doo", "Goo"},
		TodolistMessage: "Initial Todo 3",
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
			target := fmt.Sprintf(todolistByIdPath, initialDataOne.Id)
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

			for _, todolist := range todolistDB.Todolists {
				// Ensure the selected todolist is updated
				if todolist.Id == initialDataOne.Id {
					assert.Equal(t, len(test.Tags), len(resBody.Data.Tags))
					assert.Equal(t, len(test.Tags), len(todolist.Tags))

					assert.Equal(t, test.Tags, resBody.Data.Tags)
					assert.Equal(t, test.Tags, todolist.Tags)

					assert.Equal(t, test.TodolistMessage, resBody.Data.TodolistMessage)
					assert.Equal(t, test.TodolistMessage, todolist.TodolistMessage)

					assert.Equal(t, resBody.Data.CreatedAt, todolist.CreatedAt)
					assert.Equal(t, resBody.Data.UpdatedAt, todolist.UpdatedAt)

					assert.Equal(t, initialDataOne.Id, resBody.Data.Id)

					if test.Done == false {
						assert.False(t, resBody.Data.Done)
						assert.False(t, todolist.Done)
					} else {
						assert.True(t, resBody.Data.Done)
						assert.True(t, todolist.Done)
					}

					break
				}
			}

			// Ensure another todolist is not updated
			for _, todolist := range todolistDB.Todolists {
				var anotherTodolist domain.Todolist

				if todolist.Id == initialDataTwo.Id {
					anotherTodolist = initialDataTwo
				} else if todolist.Id == initialDataThree.Id {
					anotherTodolist = initialDataThree
				} else {
					continue
				}

				assert.Equal(t, anotherTodolist.Id, todolist.Id)
				assert.Equal(t, anotherTodolist.Done, todolist.Done)
				assert.ElementsMatch(t, anotherTodolist.Tags, todolist.Tags)
				assert.Equal(t, anotherTodolist.TodolistMessage, todolist.TodolistMessage)
				assert.Equal(t, anotherTodolist.CreatedAt, todolist.CreatedAt)
				assert.Equal(t, anotherTodolist.UpdatedAt, todolist.UpdatedAt)
			}
		})
	}

	resetTodolistsDB()
}

func TestUpdateTodolistFailed(t *testing.T) {
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
				target = fmt.Sprintf(todolistByIdPath, test.Id)
			} else {
				target = fmt.Sprintf(todolistByIdPath, initialData.Id)
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

			resBody := &web.WebResponse[string]{}

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
	initialDataOne := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Foo"},
		TodolistMessage: "Initial Todo 1",
	})

	time.Sleep(300 * time.Millisecond)

	initialDataTwo := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Boo"},
		TodolistMessage: "Initial Todo 2",
	})

	time.Sleep(300 * time.Millisecond)

	initialDataThree := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
		Tags:            []string{"Doo"},
		TodolistMessage: "Initial Todo 3",
	})

	time.Sleep(300 * time.Millisecond)

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
	assert.ElementsMatch(t, []domain.Todolist{initialDataTwo, initialDataThree}, todolistDB.Todolists)

	resetTodolistsDB()
}

func TestDeleteTodolistFailed(t *testing.T) {
	initialData := repository.NewTodolistRepositoryImpl().Save(dbPath, domain.Todolist{
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
	assert.ElementsMatch(t, []domain.Todolist{initialData}, todolistDB.Todolists)

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

		resBody := &web.WebResponse[[]domain.Todolist]{}

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
