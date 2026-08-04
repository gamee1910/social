package routes

import (
	"errors"
	"net/http"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/utils"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	_ = utils.ResponseJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	_ = utils.ResponseJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	_ = utils.ResponseJSONError(w, http.StatusNotFound, err.Error())
}

func ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	_ = utils.ResponseJSONError(w, http.StatusConflict, err.Error())
}

func ResponseValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	_ = utils.WriteJSON(w, http.StatusBadRequest, err)
}

func HandleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	var notFoundErr *domain.NotFoundError
	var conflictErr *domain.ConflictError

	switch {
	case errors.As(err, &notFoundErr):
		NotFoundError(w, r, err)
	case errors.As(err, &conflictErr):
		ConflictError(w, r, err)
	default:
		InternalServerError(w, r, err)
	}
}
