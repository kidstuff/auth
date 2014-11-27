package auth

import (
	"errors"
)

var (
	ErrInvalidCredential = errors.New("kidstuff/auth: Invalid emaill or password")
	ErrInvalidId         = errors.New("kidstuff/auth: Invalif Id")
	ErrInvalidActiveCode = errors.New("kidstuff/auth: Invalid activate code")
	ErrInvalidResetCode  = errors.New("kidstuff/auth: Invalid reset code")
	ErrForbidden         = errors.New("kidstuff/auth: Forbidden")
	ErrPwdMismatch       = errors.New("kidstuff/auth: Pwd and PwdRepeat doesn't match")
	ErrNoKeyProvided     = errors.New("kidstuff/auth: no key provided")
	ErrMailFailed        = errors.New("kidstuff/auth: email send failed")
)
