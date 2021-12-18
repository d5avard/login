package controllers

import (
	"html/template"
	"log"
	"path"

	"github.com/d5avard/diary/utils"
	"github.com/go-playground/validator"
)

var (
	Tmpl     *template.Template
	Validate *validator.Validate
)

const (
	templates string = "./templates/*/*"
)

func init() {
	wd := utils.GetWD()
	p := path.Join(wd, templates)
	log.Println("path templates:", p)
	Tmpl = template.Must(template.ParseGlob(p))

	Validate = validator.New()
}
