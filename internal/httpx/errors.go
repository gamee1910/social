package httpx

import (
	"log"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = ResponseJSONError(w, http.StatusInternalServerError, "the server encountered a problme")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad r error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = ResponseJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = ResponseJSONError(w, http.StatusNotFound, "resources not found")
}

func ResponseValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	log.Printf("w validation error: [%s] - path: [%s] - error: [%s]", r.Method, r.URL.Path, err)
	_ = WriteJSON(w, http.StatusBadRequest, err)
}
