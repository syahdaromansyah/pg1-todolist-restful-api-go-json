package web

type TodolistUpdateRequest struct {
	Id              string   `validate:"required" json:"id"`
	Done            bool     `json:"done"`
	Tags            []string `validate:"min=0,dive,required" json:"tags"`
	TodolistMessage string   `validate:"required,min=1,max=100" json:"todolistMessage"`
}
