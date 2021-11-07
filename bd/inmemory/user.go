package inmemory

import (
	"log"

	"github.com/d5avard/login/models"
)

var UserDB map[string]*models.User

func init() {
	UserDB = map[string]*models.User{}
}

func AddUser(user *models.User) {
	UserDB[user.UUID] = user
	log.Println("create user:", user.UUID)
}

func FindUser(user string) *models.User {
	for _, u := range UserDB {
		if u.Email == user {
			return u
		}
	}
	return nil
}
