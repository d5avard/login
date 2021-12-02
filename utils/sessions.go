package utils

import (
	"errors"
	"log"
	"net/http"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
	uuid "github.com/satori/go.uuid"
)

const Session string = "session"

// create user inmemory
// create a cookie (id, userid)
func CreateUser(w http.ResponseWriter, u *models.User) error {
	inmemory.AddUser(u)
	// TOTO: Add user cookie
	// id := uuid.NewV4().String()
	// c := http.Cookie{
	// 	Name:  id,
	// 	Value: u.UUID,
	// }
	// http.SetCookie(w, &c)

	return nil
}

func CreateSession(w http.ResponseWriter, u *models.User) {
	// create session in memory
	sessionId := uuid.NewV4().String()
	inmemory.AddSession(sessionId, u.UUID)

	// write a cookie session
	cs := http.Cookie{
		Name:   Session,
		Value:  u.UUID,
		MaxAge: 60 * 30,
	}
	http.SetCookie(w, &cs)
}

func FindSesionCookie(r *http.Request) (*http.Cookie, error) {
	var err error
	var c *http.Cookie
	c, err = r.Cookie(Session)
	if err != nil {
		log.Println("err, unable to find session cookie")
		return nil, err
	}
	return c, nil
}

// check if cookie exists ("session, userid")
// extract userid
// ckeck if sessions exists in memerory
// check if user exist in memory
func AlreadySignedIn(r *http.Request) error {
	var err error
	var c *http.Cookie

	c, err = r.Cookie(Session)
	if err != nil {
		err = errors.New("session cookie does not exist")
		log.Println(err.Error())
		return err
	}
	userid := c.Value

	if sessionId := inmemory.FindSession(userid); sessionId == "" {
		err = errors.New("session in memory does not exist")
		log.Println(err.Error())
		return err
	}

	if user := inmemory.GetUserById(userid); user == nil {
		err = errors.New("user in memory does not exist")
		log.Println(err.Error())
		return err
	}

	return nil
}
