package scheme

import "github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"

type TodolistDB struct {
	Todolists []domain.Todolist `json:"todolists"`
	Total     uint              `json:"total"`
}
