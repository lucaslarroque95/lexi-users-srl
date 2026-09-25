package models

type User struct {
	ID       string
	Email    string
	Password string
	Roles    []Role
}
