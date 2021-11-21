package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (re *RestError) Write(w http.ResponseWriter) {
	SetContentType(w, ContentTypeJSON)
	w.WriteHeader(re.Status)
	if err := json.NewEncoder(w).Encode(re); err != nil {
		log.Fatal("cannot json encode resterror")
	}
}
