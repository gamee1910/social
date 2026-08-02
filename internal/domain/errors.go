package domain

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrVersionConflict = errors.New("version conflict")
)
