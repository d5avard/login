package utils

import (
	"log"
	"os"
)

var Service string
var Revision string
var Port string
var ProjectId string

func init() {
	Service = os.Getenv("K_SERVICE")
	Revision = os.Getenv("K_REVISION")
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	if ProjectId = os.Getenv("GOOGLE_APPLICATION_PROJECT_ID"); ProjectId == "" {
		log.Fatalln("no project id is configured")
	}
}
