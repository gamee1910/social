package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
func FormatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)

	var validationErrors validator.ValidationErrors

	ok := errors.As(err, &validationErrors)

	if !ok {
		return errs
	}
	for _, fieldError := range validationErrors {
		param := fieldError.Param()
		field := fieldError.Field()

		switch fieldError.Tag() {
		case "required":
			errs[field] = fmt.Sprintf("%s is required", field)

		case "email":
			errs[field] = fmt.Sprintf("%s must be a valid email", field)

		case "min":
			errs[field] = fmt.Sprintf("%s must be at least %s characters", field, param)

		case "max":
			errs[field] = fmt.Sprintf("%s must not exceed %s characters", field, param)

		default:
			errs[field] = fmt.Sprintf("%s is invalid", field)
		}
	}

	return errs
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1MegaByte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(data)
}

func ResponseJSONError(w http.ResponseWriter, status int, message string) error {
	type envelop struct {
		Error string `json:"error"`
	}
	return WriteJSON(w, status, &envelop{Error: message})
}
