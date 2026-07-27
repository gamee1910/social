package httpx

import (
	"log"
	"net/http"
)

func InternalServerError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("internal server error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = ResponseJSONError(response, http.StatusInternalServerError, "the server encountered a problme")
}

func BadRequestError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("bad request error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = ResponseJSONError(response, http.StatusBadRequest, err.Error())
}

func NotFoundError(response http.ResponseWriter, request *http.Request, err error) {
	log.Printf("not found error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = ResponseJSONError(response, http.StatusNotFound, "resources not found")
}

func ResponseValidationError(response http.ResponseWriter, request *http.Request, err map[string]string) {
	log.Printf("response validation error: [%s] - path: [%s] - error: [%s]", request.Method, request.URL.Path, err)
	_ = WriteJSON(response, http.StatusBadRequest, err)
}
