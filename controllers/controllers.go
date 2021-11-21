package controllers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/utils"
	"github.com/go-playground/validator"
)

var (
	Tmpl     *template.Template
	Validate *validator.Validate
)

const (
	templates string = "./templates/*/*"
	session   string = "session"
)

func init() {
	wd := utils.GetWD()
	p := path.Join(wd, templates)
	log.Println("path templates:", p)
	Tmpl = template.Must(template.ParseGlob(p))

	Validate = validator.New()
}

// check if cookie exists
// extract user id
// ckeck if sessions exists
// check if user exist
func AlreadySignedIn(r *http.Request) error {
	var ok bool
	var err error
	var c *http.Cookie

	c, err = r.Cookie(session)
	if err != nil {
		err = errors.New("session cookie does not exist")
		log.Println(err.Error())
		return err
	}
	id := c.Value

	if _, ok = inmemory.SessionDB[id]; !ok {
		err = errors.New("session in memory does not exist")
		log.Println(err.Error())
		return err
	}

	_, err = r.Cookie(id)
	if err != nil {
		err = errors.New("user cookie does not exist")
		log.Println(err.Error())
		return err
	}

	if _, ok = inmemory.UserDB[id]; !ok {
		err = errors.New("user in memory does not exist")
		log.Println(err.Error())
		return err
	}

	return nil
}
