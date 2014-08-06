// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/base64"
	"errors"
	"github.com/gorilla/securecookie"
	"time"
)

var (
	ErrInvalidId       = errors.New("auth: invalid id")
	ErrInvalidEmail    = errors.New("auth: invalid email address")
	ErrDuplicateEmail  = errors.New("auth: duplicate email address")
	ErrInvalidPassword = errors.New("auth: invalid password")
	ErrNotLogged       = errors.New("auth: no login user found")
)

type User struct {
	Id           interface{}       `bson:"_id"`
	Email        *string           `bson:"Email" json:",omitempty"`
	OldPwd       []Password        `bson:"OldPwd" json:"-"`
	Pwd          *Password         `bson:"Pwd" json:"-"`
	LastActivity *time.Time        `bson:"LastActivity"  json:",omitempty"`
	Info         *UserInfo         `bson:"Info" json:",omitempty"`
	Privilege    []string          `bson:"Privilege" json:",omitempty"`
	Approved     *bool             `bson:"Approved" json:",omitempty"`
	ConfirmCodes map[string]string `bson:"ConfirmCodes" json:"-"`
	BriefGroups  []BriefGroup      `bson:"BriefGroups" json:",omitempty"`
}

func (u *User) SetPassword(pwd string) error {
	panic("")
}

func (u *User) ChangePassword(pwd string) error {
	var err error
	u.OldPwd = append(u.OldPwd, *u.Pwd)
	hpwd, err := HashPwd(pwd)
	if err != nil {
		return err
	}

	u.Pwd = &hpwd
	return nil
}

func (u *User) ComparePassword(pwd string) error {
	pwdBytes := []byte(pwd)
	tmp := make([]byte, len(pwdBytes)+len(u.Pwd.Salt))
	copy(tmp, pwdBytes)
	tmp = append(tmp, u.Pwd.Salt...)
	if err := bcrypt.CompareHashAndPassword(u.Pwd.Hashed, tmp); err != nil {
		return err
	}

	return nil
}

// ValidConfirmCode valid the code for specific key of the user specify by id.
// Re-generate or delete code for that key if need.
func (u *User) ValidConfirmCode(key, code string, regen, del bool) bool {
	if u.ConfirmCodes[key] == code {
		if del {
			delete(u.ConfirmCodes, key)
		}

		if regen {
			u.ConfirmCodes[key] = base64.URLEncoding.EncodeToString(securecookie.
				GenerateRandomKey(64))
		}

		return true
	}

	return false
}

type Password struct {
	Hashed []byte    `bson:"Hashed"`
	Salt   []byte    `bson:"Salt"`
	InitAt time.Time `bson:"InitAt"`
}

type UserInfo struct {
	FirstName  *string    `bson:"FirstName" json:",omitempty"`
	LastName   *string    `bson:"LastName" json:",omitempty"`
	MiddleName *string    `bson:"MiddleName" json:",omitempty"`
	NickName   *string    `bson:"NickName" json:",omitempty"`
	BirthDay   *time.Time `bson:"BirthDay" json:",omitempty"`
	JoinDay    *time.Time `bson:"JoinDay" json:",omitempty"`
	Address    []Address  `bson:"Address" json:",omitempty"`
	Phone      []string   `bson:"Phone" json:",omitempty"`
}

type Address struct {
	Country  *string `bson:"Country" json:",omitempty"`
	State    *string `bson:"State" json:",omitempty"`
	City     *string `bson:"City" json:",omitempty"`
	District *string `bson:"District" json:",omitempty"`
	Street   *string `bson:"Street" json:",omitempty"`
}

type UserManager interface {
	// Add adds an user to database with email and password;
	// If app is false, the user is waiting to be approved.
	// Implement of this method should valid email, pwd and make sure the user
	// email are unique.
	// It returns an error describes the first issue encountered, if any.
	Add(email, pwd string, app bool) (*User, error)
	// AddDetail add a User with full detail to database.
	// Implement of this method should valid email, pwd and make sure the user
	// email are unique.
	// It returns an error describes the first issue encountered, if any.
	AddDetail(*User) (*User, error)
	// UpdateDetail changes detail of the User.
	// It returns an error describes the first issue encountered, if any.
	UpdateDetail(*User) error
	// Delete deletes an user from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	Delete(id interface{}) error
	// Find finds the user with the given id;
	// Its returns an ErrNotFound if the user's id was not found.
	Find(id interface{}) (*User, error)
	// FindByEmail like Find but receive an email
	FindByEmail(email string) (*User, error)
	// FindAll finds and return a slice of group.
	// If limit < 0 the mean using the default upper limit.
	// If limit == 0 return empty result with error indicate no result found.
	// If limit can't be greater than the default upper limit.
	// Specific fields name for porjection select.
	FindAll(limit int, offsetId interface{}, fields []string) ([]*User, error)
	// FindAllOnline finds and return a slice of current Loged user.
	// See FindAll for the usage.
	FindAllOnline(limit int, offsetId interface{}, fields []string) ([]*User, error)
	// Get gets the infomations and update the LastActivity of the current
	// loged user by the token (given by Login method);
	// It returns an error describes the first issue encountered, if any.
	Get(token string) (*User, error)
	// Login logs user in by given user id.
	// Stay is the duration to keep the user Login state.
	// It returns a token string, use the token to keep track on the user with
	// Get or Logout.
	Login(id interface{}, stay time.Duration) (string, error)
	// Logout logs the current user out.
	Logout(token string) error
}
