package models

//Car models...
type Car struct {
	Mark   string `json:"mark"`
	MaxSpeed  int `json:"max_speed"`
	Distance int `json:"distance"`
	Handler string `json:"handler"`
	Stock string `json:"stock"`
}
