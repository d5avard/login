package controllers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// check if the user is already signedin
	if !AlreadySignedIn(r) {
		// if not sign, redirect to signin
		log.Println("error, user is not signin")
		w.Header().Set("Location", "/signin")
		w.WriteHeader(http.StatusSeeOther)
	}

	Tmpl.ExecuteTemplate(w, "index", nil)
}
