package utils

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func SetContentType(w http.ResponseWriter, ct string) {
	w.Header().Add(HeaderContentType, ct)
}

func NotFound(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
