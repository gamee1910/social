package domain

import (
	"errors"
	"fmt"
)

// Sentinel errors — dùng để so sánh qua errors.Is
var (
	ErrNotFound        = errors.New("resource not found")
	ErrVersionConflict = errors.New("version conflict")
)

// NotFoundError — typed error cho HTTP mapping (thêm context resource nào bị thiếu)
type NotFoundError struct {
	Resource string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Resource)
}

// ConflictError — typed error cho optimistic locking conflict
type ConflictError struct {
	Resource string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s has been modified by another request, please try again", e.Resource)
}
