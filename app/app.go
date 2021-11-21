package app

import (
	"log"
	"net/http"

	"github.com/d5avard/login/controllers"
	"github.com/julienschmidt/httprouter"
)

func App() {
	router := httprouter.New()
	controllers.Routes(router)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println("err:", err.Error())
	}
}
