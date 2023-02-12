package web

type TodolistResponse struct {
	Id              string   `json:"id"`
	Done            bool     `json:"done"`
	Tags            []string `json:"tags"`
	TodolistMessage string   `json:"todolistMessage"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}
