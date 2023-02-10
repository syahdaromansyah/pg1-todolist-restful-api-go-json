package lib

import (
	"github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/helper"

	"github.com/jaevor/go-nanoid"
)

func GetRandomStdId32() string {
	idStd32, err := nanoid.Standard(32)
	helper.DoPanicIfError(err)

	return idStd32()
}
