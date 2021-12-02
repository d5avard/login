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
	if err = utils.AlreadySignedIn(r); err != nil {
		re := utils.NewUnauthorized(err)
		re.Write(w)
		return
	}
	if err = Tmpl.ExecuteTemplate(w, "index", nil); err != nil {
		log.Println("error:", err.Error())
	}
}
