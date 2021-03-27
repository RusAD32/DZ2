package models

//User model ...
type User struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}
