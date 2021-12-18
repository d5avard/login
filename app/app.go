package app

import (
	"log"
	"net/http"

	"github.com/d5avard/diary/controllers"
	"github.com/d5avard/diary/utils"
	"github.com/julienschmidt/httprouter"
)

func App() {
	router := httprouter.New()
	controllers.Routes(router)

	log.Printf("listening on port %s", utils.Port)
	if err := http.ListenAndServe(":"+utils.Port, router); err != nil {
		log.Println("err:", err.Error())
	}
}
