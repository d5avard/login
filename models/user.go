package models

type User struct {
	UUID         string
	Email        string
	HashPassword []byte
}
