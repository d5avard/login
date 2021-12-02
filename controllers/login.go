package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/d5avard/login/bd/inmemory"
	"github.com/d5avard/login/models"
	"github.com/d5avard/login/utils"
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupForm struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,eqfield=PasswordConfirm,max=32"`
	PasswordConfirm string `json:"passwordconfirm" validate:"required,max=32"`
}

type SigninForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=32"`
	Loggedin bool   `json:"loggedin"`
}

func Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	if err = utils.AlreadySignedIn(r); err == nil {
		utils.SetLocation(w, "/index")
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		Tmpl.ExecuteTemplate(w, "signin", nil)
	} else if r.Method == "POST" {
		form := SigninForm{}
		if err = json.NewDecoder(r.Body).Decode(&form); err != nil {
			log.Println("error:", err)
			re := utils.NewStatusBadRequest(err)
			re.Write(w)
			return
		}

		if err = Validate.Struct(form); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("error field:", e.Field())
				fmt.Println("error:", e)
			}
			re := utils.NewUnauthorized(err)
			re.Write(w)
			return
		}

		// find the user
		var u *models.User
		if u = inmemory.FindUser(form.Email); u == nil {
			// error, cannot find the user
			re := utils.NewUnauthorized(errors.New("cannot find the user"))
			re.Write(w)
			return
		}

		// validate the password
		if err := bcrypt.CompareHashAndPassword(u.HashPassword, []byte(form.Password)); err != nil {
			// error, the password is not valid
			re := utils.NewUnauthorized(err)
			re.Write(w)
			return
		}

		utils.CreateSession(w, u)

		log.Println("user signedin:", u.UUID)

		utils.SetLocation(w, "/index")
		w.WriteHeader(http.StatusOK)
	}
}

func Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	if err = utils.AlreadySignedIn(r); err == nil {
		utils.SetLocation(w, "/index")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		Tmpl.ExecuteTemplate(w, "signup", nil)
	} else if r.Method == "POST" {
		var u *models.User
		form := SignupForm{}
		if err = json.NewDecoder(r.Body).Decode(&form); err != nil {
			log.Println("error:", err)
			re := utils.NewInternalServerError(err)
			re.Write(w)
			return
		}

		if err = Validate.Struct(form); err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("error field:", e.Field())
				fmt.Println("error:", e)
			}
			re := utils.NewUnauthorized(err)
			re.Write(w)
			return
		}

		if u = inmemory.FindUser(form.Email); u != nil {
			// Error, the user already exist
			re := utils.NewUnauthorized(err)
			re.Write(w)
			return
		}

		// create a session cookie  (session, userid)
		// create a user cookie (id, userid)
		// add session in memory (sessionid, userid)
		// add user in memory (userid, jsonstring)
		u = &models.User{}
		u.UUID = uuid.NewV4().String()
		u.Email = form.Email
		var hash []byte
		if hash, err = bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.MinCost); err != nil {
			re := utils.NewInternalServerError(err)
			re.Write(w)
			return
		}
		u.HashPassword = hash
		utils.CreateUser(w, u)

		log.Println("user signedup:", u.UUID)

		utils.SetLocation(w, "/signin")
		w.WriteHeader(http.StatusOK)
	}
}

func Signout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var c *http.Cookie
	if err = utils.AlreadySignedIn(r); err != nil {
		re := utils.NewUnauthorized(err)
		re.Write(w)
		return
	}

	if c, err = utils.FindSesionCookie(r); err != nil {
		re := utils.NewInternalServerError(err)
		re.Write(w)
	}

	inmemory.DeleteSession(c.Value)
	log.Println("delete session", c.Value)

	id := c.Value
	c.Value = ""
	c.MaxAge = -1
	http.SetCookie(w, c)
	log.Println("delete session cookie", id)

	utils.SetLocation(w, "/signin")
	w.WriteHeader(http.StatusOK)
}
