package helper

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/exception"
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/web"
)

func internalServerError(writer http.ResponseWriter, request *http.Request, err any) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	webResponse := web.WebResponse[string]{
		Code:   http.StatusInternalServerError,
		Status: "failed",
		Data:   "something went wrong",
	}

	WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err any) bool {
	exception, isError := err.(exception.NotFoundError)

	if isError {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse[string]{
			Code:   http.StatusNotFound,
			Status: "failed",
			Data:   exception.Error,
		}

		WriteToResponseBody(writer, webResponse, http.StatusNotFound)

		return true
	} else {
		return false
	}
}

func reqBodyMalformedError(writer http.ResponseWriter, request *http.Request, err any) bool {
	_, isError := err.(exception.ReqBodyMalformedError)

	if isError {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse[string]{
			Code:   http.StatusBadRequest,
			Status: "failed",
			Data:   "request body is invalid",
		}

		WriteToResponseBody(writer, webResponse, http.StatusBadRequest)

		return true
	} else {
		return false
	}
}

func validationError(writer http.ResponseWriter, request *http.Request, err any) bool {
	_, isError := err.(validator.ValidationErrors)

	if isError {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse[string]{
			Code:   http.StatusBadRequest,
			Status: "failed",
			Data:   "request body is invalid",
		}

		WriteToResponseBody(writer, webResponse, http.StatusBadRequest)

		return true
	} else {
		return false
	}
}

func HttpRouterPanicHandler(writer http.ResponseWriter, request *http.Request, err any) {
	if notFoundError(writer, request, err) {
		return
	}

	if reqBodyMalformedError(writer, request, err) {
		return
	}

	if validationError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}
