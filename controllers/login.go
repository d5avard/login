package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type formSignup struct {
	email           string `validate:"required,email"`
	password        string `validate:"required,eqfield=ConfirmPassword"`
	confirmPassword string `validate:"required"`
}

type formSignin struct {
	email    string `validate:"required,email"`
	password string `validate:"required"`
}

func Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if AlreadySignedIn(r) {
		w.Header().Set("Location", "/index")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Tmpl.ExecuteTemplate(w, "signin", nil)
	} else if r.Method == "POST" {
		var err error

		// get parameters
		r.ParseForm()
		form := formSignin{}
		form.email = r.FormValue("email")
		form.password = r.FormValue("password")

		if err = Validate.Struct(form); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("Error field:", e.Field())
			}

			Tmpl.ExecuteTemplate(w, "signup", nil)
			return
		}

		// find the user
		var u *models.User
		if u = inmemory.FindUser(form.email); u == nil {
			// error, cannot find the user
			Tmpl.ExecuteTemplate(w, "signin", nil)
			return
		}

		// validate the password
		if err := bcrypt.CompareHashAndPassword(u.HashPassword, []byte(form.password)); err != nil {
			// error, the password is not valid
			Tmpl.ExecuteTemplate(w, "signin", nil)
			return
		}

		// create session in memory
		inmemory.AddSession(u.UUID, u.Email)

		// write a cookie session
		cs := http.Cookie{
			Name:   session,
			Value:  u.UUID,
			MaxAge: 60 * 30,
		}
		http.SetCookie(w, &cs)

		// redirect to /index
		w.Header().Set("Location", "/index")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if AlreadySignedIn(r) {
		w.Header().Set("Location", "/index")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Tmpl.ExecuteTemplate(w, "signup", nil)
	} else if r.Method == "POST" {
		var u *models.User
		var err error

		form := formSignup{}
		form.email = r.FormValue("email")
		form.password = r.FormValue("password")
		form.confirmPassword = r.FormValue("password-validate")

		if err = Validate.Struct(form); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("Error field:", e.Field())
			}

			Tmpl.ExecuteTemplate(w, "signup", nil)
			return
		}

		if u = inmemory.FindUser(form.email); u != nil {
			// Error, the user already exist
			Tmpl.ExecuteTemplate(w, "signup", nil)
			return
		}

		u = &models.User{}
		u.UUID = uuid.NewV4().String()
		u.Email = form.email
		var hash []byte
		if hash, err = bcrypt.GenerateFromPassword([]byte(form.password), bcrypt.MinCost); err != nil {
			http.Error(w, "err, internal server error", http.StatusInternalServerError)
			return
		}
		u.HashPassword = hash
		inmemory.AddUser(u)
		c := http.Cookie{
			Name:  u.UUID,
			Value: "true",
		}
		http.SetCookie(w, &c)
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func Signout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !AlreadySignedIn(r) {
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	cs, err := r.Cookie(session)
	if err != nil {
		log.Println("err, unable to find session cookie")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	inmemory.DeleteSession(cs.Value)
	log.Println("delete session", cs.Value)

	id := cs.Value
	cs.Value = ""
	cs.MaxAge = -1
	http.SetCookie(w, cs)
	log.Println("delete session cookie", id)

	w.Header().Set("Location", "/signin")
	w.WriteHeader(http.StatusSeeOther)
}
