package controllers

import (
	"encoding/json"
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
	Password        string `json:"password" validate:"required,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"passwordconfirm" validate:"required"`
}

type SigninForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Loggedin bool   `json:"loggedin"`
}

func Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	if err = AlreadySignedIn(r); err == nil {
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

		// find the user
		var u *models.User
		if u = inmemory.FindUser(form.Email); u == nil {
			// error, cannot find the user
			re := utils.NewUnauthorized(err)
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

		// create session in memory
		inmemory.AddSession(u.UUID, u.Email)

		// write a cookie session
		cs := http.Cookie{
			Name:   session,
			Value:  u.UUID,
			MaxAge: 60 * 30,
		}
		http.SetCookie(w, &cs)

		log.Println("user signedin:", u.UUID)

		utils.SetLocation(w, "/index")
		w.WriteHeader(http.StatusOK)
	}
}

func Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	if err = AlreadySignedIn(r); err == nil {
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
	if err := AlreadySignedIn(r); err != nil {
		re := utils.NewUnauthorized(err)
		re.Write(w)
		return
	}

	cs, err := r.Cookie(session)
	if err != nil {
		log.Println("err, unable to find session cookie")
		re := utils.NewInternalServerError(err)
		re.Write(w)
		return
	}

	inmemory.DeleteSession(cs.Value)
	log.Println("delete session", cs.Value)

	id := cs.Value
	cs.Value = ""
	cs.MaxAge = -1
	http.SetCookie(w, cs)
	log.Println("delete session cookie", id)

	utils.SetLocation(w, "/signin")
	w.WriteHeader(http.StatusOK)
}
