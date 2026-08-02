package httputil

import (
	"errors"
	"fmt"

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
