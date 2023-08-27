package repository

import "github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"

type TodolistRepository interface {
	Save(dbPath string, todolistRequest domain.Todolist) domain.Todolist
	Update(dbPath string, todolistRequest domain.Todolist) (domain.Todolist, error)
	Delete(dbPath string, todolistIdParam string) error
	FindAll(dbPath string) []domain.Todolist
}
