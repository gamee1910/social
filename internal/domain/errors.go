package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("resource not found")
	ErrVersionConflict = errors.New("version conflict")
)

type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

type ConflictError struct {
	Resource string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s has been modified by another request, please try again", e.Resource)
}
