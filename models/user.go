package models

type User struct {
	UUID         string
	Email        string
	HashPassword []byte
}

func NewUser() *User {
	u := User{}
	return &u
}
