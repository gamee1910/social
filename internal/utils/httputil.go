package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

// ReadJSON decodes a JSON request body into data.
// w is required by http.MaxBytesReader to enforce the body size limit.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1 << 20 // 1 MiB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
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
	slog.ErrorContext(r.Context(), "internal server error",
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
	)
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
	switch {
	case errors.Is(err, domain.ErrNotFound):
		NotFoundError(w, r, err)
	case errors.Is(err, domain.ErrVersionConflict):
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
