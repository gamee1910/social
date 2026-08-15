package domain

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrVersionConflict       = errors.New("version conflict")
	ErrUnauthorized          = errors.New("email or password not correct")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
)
