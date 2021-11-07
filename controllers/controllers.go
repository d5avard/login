package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/go-playground/validator"
)

var (
	Tmpl     *template.Template
	Validate *validator.Validate
)

const (
	templates string = "./templates/*"
	session   string = "session"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("err:", err.Error())
	}
	p := path.Join(cwd, templates)
	log.Println("path templates:", p)
	Tmpl = template.Must(template.ParseGlob(p))

	Validate = validator.New()
}

// check if cookie exists
// extract user id
// ckeck if sessions exists
// check if user exist
func AlreadySignedIn(r *http.Request) bool {
	var ok bool
	c, err := r.Cookie(session)
	if err != nil {
		return false
	}

	if _, ok = inmemory.SessionDB[c.Value]; !ok {
		return false
	}

	if _, ok = inmemory.UserDB[c.Value]; !ok {
		return false
	}

	return true
}
