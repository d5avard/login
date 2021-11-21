package controllers

import (
	"log"
	"net/http"

	"github.com/d5avard/login/utils"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// check if the user is already signedin
	var err error
	if err = AlreadySignedIn(r); err != nil {
		// if not sign, redirect to signin
		re := utils.RestError{Message: err.Error(), Status: http.StatusUnauthorized, Error: "unauthorized"}
		re.Write(w)
		return
	}
	if err = Tmpl.ExecuteTemplate(w, "index", nil); err != nil {
		log.Println("error:", err.Error())
	}
}
