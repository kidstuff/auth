package auth

import (
	"errors"
)

var (
	ErrInvalidCredential = errors.New("kidstuff/auth: Invalid emaill or password")
	ErrInvalidId         = errors.New("kidstuff/auth: Invalif Id")
	ErrInvalidActiveCode = errors.New("kidstuff/auth: Invalid activate code")
	ErrForbidden         = errors.New("kidstuff/auth: Forbidden")
)
