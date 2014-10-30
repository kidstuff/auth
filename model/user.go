// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	// "code.google.com/p/go.crypto/bcrypt"
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
	ErrDuplicateName   = errors.New("auth: duplicate Group Name")
)

type User struct {
	Id           *string           `bson:"Id,omitempty"`
	Email        *string           `bson:"Email" json:",omitempty"`
	Pwd          *Password         `bson:"Pwd" json:"-"`
	LastActivity *time.Time        `bson:"LastActivity"  json:",omitempty"`
	Privileges   []string          `bson:"Privileges" json:",omitempty"`
	Approved     *bool             `bson:"Approved" json:",omitempty"`
	ConfirmCodes map[string]string `bson:"ConfirmCodes" json:"-"`
	Profile      *Profile          `bson:"Profile,omitempty" json:",omitempty"`
	Groups       []*Group          `bson:"Groups,omitempty" json:",omitempty"`
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

type Profile struct {
	FirstName  *string    `bson:"FirstName" json:",omitempty"`
	LastName   *string    `bson:"LastName" json:",omitempty"`
	MiddleName *string    `bson:"MiddleName" json:",omitempty"`
	NickName   *string    `bson:"NickName" json:",omitempty"`
	BirthDay   *time.Time `bson:"BirthDay" json:",omitempty"`
	JoinDay    *time.Time `bson:"JoinDay" json:",omitempty"`
	Addresses  []Address  `bson:"Addresses" json:",omitempty"`
	Phones     []string   `bson:"Phones" json:",omitempty"`
}

type Address struct {
	Country  *string `bson:"Country" json:",omitempty"`
	State    *string `bson:"State" json:",omitempty"`
	City     *string `bson:"City" json:",omitempty"`
	District *string `bson:"District" json:",omitempty"`
	Street   *string `bson:"Street" json:",omitempty"`
}

type Group struct {
	Id        *string    `bson:"Id,omitempty"`
	Name      *string    `bson:"Name,omitempty" json:",omitempty"`
	Info      *GroupInfo `bson:"Info,omitempty" json:",omitempty"`
	Privilege []string   `bson:"Privilege,omitempty" json:",omitempty"`
}

type GroupInfo struct {
	Description *string `bson:"Description" json:",omitempty"`
}
