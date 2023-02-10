package repository

import "github.com/syahdaromansyah/pg1-todolist-restful-api-go-json/model/domain"

type TodolistRepository interface {
	Save(dbPath string, todolistRequest domain.Todolist) domain.Todolist
	Update(dbPath string, todolistRequest domain.Todolist) domain.Todolist
	Delete(dbPath string, todolistRequest domain.Todolist)
	FindById(dbPath string, todolistIdParam string) (domain.Todolist, error)
	FindAll(dbPath string) []domain.Todolist
}
