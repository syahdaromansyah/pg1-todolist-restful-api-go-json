package helper

import (
	"encoding/json"
	"net/http"

	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/exception"
)

func ReadFromRequestBody(request *http.Request, result any) {
	jsonDecoder := json.NewDecoder(request.Body)
	err := jsonDecoder.Decode(result)

	if err != nil {
		panic(exception.NewReqBodyMalformedError(err.Error()))
	}
}

func WriteToResponseBody(writer http.ResponseWriter, response any) {
	writer.Header().Add("Content-Type", "application/json")
	jsonEncoder := json.NewEncoder(writer)
	err := jsonEncoder.Encode(response)
	DoPanicIfError(err)
}
