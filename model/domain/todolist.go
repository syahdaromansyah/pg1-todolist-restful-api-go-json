package domain

type Todolist struct {
	Id              string
	Done            bool
	Tags            []string
	TodolistMessage string
	CreatedAt       string
	UpdatedAt       string
}
