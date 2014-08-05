package model

import (
	"regexp"
)

type FormatChecker interface {
	PasswordValidate(string) bool
	EmailValidate(string) bool
}

type SimpleChecker struct {
	emailregex *regexp.Regexp
	pwdlen     int
}

func NewSimpleChecker(pwdlen int) (*SimpleChecker, error) {
	var err error

	c := SimpleChecker{}
	c.emailregex, err = regexp.Compile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	c.pwdlen = pwdlen
	return &c, err
}

func (c *SimpleChecker) PasswordValidate(pwd string) bool {
	return len(pwd) >= c.pwdlen
}

func (c *SimpleChecker) EmailValidate(email string) bool {
	return c.emailregex.MatchString(email)
}
