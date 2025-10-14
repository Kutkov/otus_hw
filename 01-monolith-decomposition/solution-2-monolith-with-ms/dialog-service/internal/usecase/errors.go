package usecase

import "errors"

var (
	ErrInvalidData  = errors.New("invalid data")
	ErrUnauthorized = errors.New("unauthorized")
)
