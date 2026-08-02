package routes

import (
	"log"
	"net/http"

	"github.com/gamee1910/social/internal/config"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = config.ResponseJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = config.ResponseJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = config.ResponseJSONError(w, http.StatusNotFound, "resource not found")
}

func ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = config.ResponseJSONError(w, http.StatusConflict, err.Error())
}

func ResponseValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	log.Printf("validation error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = config.WriteJSON(w, http.StatusBadRequest, err)
}
