package controllers

import (
	"github.com/d5avard/login/utils"
	"github.com/julienschmidt/httprouter"
)

func Routes(r *httprouter.Router) {
	r.GET("/", Index)
	r.GET("/index", Index)
	r.GET("/signin", Signin)
	r.POST("/signin", Signin)
	r.GET("/signup", Signup)
	r.POST("/signup", Signup)
	r.GET("/signout", Signout)
	r.GET("/favicon.ico", utils.NotFound)
}
