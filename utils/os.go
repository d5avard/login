package utils

import (
	"log"
	"os"
)

func GetWD() string {
	env := os.Getenv("ENV")
	var wd string
	var err error
	if env == "TEST" {
		if wd = os.Getenv("WD"); wd == "" {
			log.Fatalln("error: you need to set the env variable WD")
		}
	} else if env == "" {
		if wd, err = os.Getwd(); err != nil {
			log.Fatalln("error:", err.Error())
		}
	}
	return wd
}
