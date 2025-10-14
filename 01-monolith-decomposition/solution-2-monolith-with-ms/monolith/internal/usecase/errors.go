package usecase

import "errors"

var (
	ErrInvalidData      = errors.New("invalid data")
	ErrInvalidBirthdate = errors.New("invalid birthdate")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUnauthorized     = errors.New("unauthorized")
)
