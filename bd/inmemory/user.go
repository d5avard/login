package inmemory

import (
	"log"

	"github.com/d5avard/diary/models"
)

var userDB map[string]*models.User

func init() {
	userDB = map[string]*models.User{}
}

func AddUser(user *models.User) {
	userDB[user.UUID] = user
	log.Println("create user:", user.UUID)
}

func GetUserById(userid string) *models.User {
	return userDB[userid]
}

func FindUser(user string) *models.User {
	for _, u := range userDB {
		if u.Email == user {
			return u
		}
	}
	return nil
}
