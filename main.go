package main

import (
	"log"
	"net/http"

	"github.com/d5avard/login/controllers"
	"github.com/julienschmidt/httprouter"
)

var Router *httprouter.Router

func main() {
	Router = httprouter.New()

	Router.GET("/", controllers.Index)
	Router.GET("/index", controllers.Index)
	Router.GET("/signin", controllers.Signin)
	Router.POST("/signin", controllers.Signin)
	Router.GET("/signup", controllers.Signup)
	Router.POST("/signup", controllers.Signup)
	Router.GET("/signout", controllers.Signout)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	if err := http.ListenAndServe(":8080", Router); err != nil {
		log.Println("err:", err.Error())
	}
}
