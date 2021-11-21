package utils

import (
	"net/http"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
)

func CreateUser(w http.ResponseWriter, u *models.User) error {
	inmemory.AddUser(u)
	c := http.Cookie{
		Name:  u.UUID,
		Value: "true",
	}
	http.SetCookie(w, &c)

	return nil
}
