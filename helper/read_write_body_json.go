package helper

import (
	"encoding/json"
	"net/http"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/exception"
)

func ReadFromRequestBody[T any](request *http.Request, result T) {
	jsonDecoder := json.NewDecoder(request.Body)
	err := jsonDecoder.Decode(result)

	if err != nil {
		panic(exception.NewReqBodyMalformedError(err.Error()))
	}
}

func WriteToResponseBody[T any](writer http.ResponseWriter, response T, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	jsonEncoder := json.NewEncoder(writer)
	err := jsonEncoder.Encode(response)
	DoPanicIfError(err)
}
