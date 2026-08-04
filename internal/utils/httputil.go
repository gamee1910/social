package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/domain"
	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1 MegaByte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

func ResponseJSONError(w http.ResponseWriter, status int, message any) error {
	type envelop struct {
		Error any `json:"error"`
	}
	return writeJSON(w, status, &envelop{Error: message})
}

func ResponseJSON(w http.ResponseWriter, status int, data any) error {
	type responseJSON struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &responseJSON{Data: data})
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	_ = ResponseJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	_ = ResponseJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	_ = ResponseJSONError(w, http.StatusNotFound, err.Error())
}

func ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	_ = ResponseJSONError(w, http.StatusConflict, err.Error())
}

func ResponseValidationError(w http.ResponseWriter, r *http.Request, err map[string]string) {
	_ = ResponseJSONError(w, http.StatusBadRequest, err)
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
func GetIDFromParameter(value string, r *http.Request) (int64, error) {
	id := chi.URLParam(r, value)

	valueFromParam, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return valueFromParam, nil
}
