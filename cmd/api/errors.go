package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("internal server error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = responseJSONError(response, http.StatusInternalServerError, "the server encountered a problme")
}

func (app *application) badRequestError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("bad request error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = responseJSONError(response, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("not found error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = responseJSONError(response, http.StatusNotFound, "resources not found")
}

func (app *application) responseValidationError(response http.ResponseWriter, request *http.Request, err map[string]string) {
	log.Printf("response validation error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = writeJSON(response, http.StatusBadRequest, err)
}
