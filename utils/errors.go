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

func NewInternalServerError(err error) *RestError {
	return &RestError{Message: err.Error(), Status: http.StatusInternalServerError, Error: "internal_server_error"}
}

func NewUnauthorized(err error) *RestError {
	return &RestError{Message: err.Error(), Status: http.StatusUnauthorized, Error: "unauthorized"}
}

func NewStatusBadRequest(err error) *RestError {
	return &RestError{Message: err.Error(), Status: http.StatusBadRequest, Error: "bad_request"}
}

func (re *RestError) Write(w http.ResponseWriter) {
	SetContentType(w, ContentTypeJSON)
	w.WriteHeader(re.Status)
	if err := json.NewEncoder(w).Encode(re); err != nil {
		log.Fatal("cannot json encode resterror")
	}
}
