package model

// struct used to maintain an individual todo item
type ToDo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
