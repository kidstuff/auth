// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authmodel

import (
	"regexp"
)

// FormatChecker is a helper interface, a "manager" donnot need to implement it.
type FormatChecker interface {
	// PasswordValidate validate password strength.
	PasswordValidate(string) bool
	// EmailValidate validate email format.
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

// PasswordValidate validate pwd base on the length
func (c *SimpleChecker) PasswordValidate(pwd string) bool {
	return len(pwd) >= c.pwdlen
}

// EmailValidate validate email base on a simple regex
func (c *SimpleChecker) EmailValidate(email string) bool {
	return c.emailregex.MatchString(email)
}
