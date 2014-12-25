// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authmodel

import (
	"time"
)

type Manager interface {
	// AddUser adds an user to database with email and password;
	// If app is false, the user is waiting to be approved.
	// Implement of this method should valid email, pwd and make sure the user
	// email are unique then initial the LastActivity and JoinDay.
	// Remember to generate ConfirmCodes["activate"] a secure random value.
	// It returns an error describes the first issue encountered, if any.
	AddUser(email, pwd string, app bool) (*User, error)
	// AddUserDetail add a User with full detail to database.
	// Implement of this method should valid email, pwd and make sure the user
	// email are unique then initial the LastActivity and JoinDay.
	// Remember to generate ConfirmCodes["activate"] a secure random value.
	// It returns an error describes the first issue encountered, if any.
	AddUserDetail(email, pwd string, app bool, pri []string,
		code map[string]string, profile *Profile, groupIds []string) (*User, error)
	// UpdateUserDetail changes detail of the User.
	// If a agrument is nil, the implement of tis function should not affect the
	// value of that field.
	// It returns an error describes the first issue encountered, if any.
	UpdateUserDetail(id string, pwd *string, app *bool, pri []string,
		code map[string]string, profile *Profile, groupIds []string) error
	// DeleteUser deletes an user from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	DeleteUser(id string) error
	// FindUser finds the user with the given id;
	// Its returns an ErrNotFound if the user's id was not found.
	FindUser(id string) (*User, error)
	// FindUserByEmail like Find but receive an email
	FindUserByEmail(email string) (*User, error)
	// FindAllUser finds and return a slice of users belong to some groups if specific.
	// If limit < 0 the mean using the default upper limit.
	// If limit == 0 return empty result with error indicate no result found.
	// If limit can't be greater than the default upper limit.
	// Specific fields name for porjection select.
	FindAllUser(limit int, offsetId string, fields []string, groupIds []string) ([]*User, error)
	// GetUser gets the infomations and update the LastActivity of the current
	// loged user by the token (given by Login method);
	// It returns an error describes the first issue encountered, if any.
	GetUser(token string) (*User, error)
	// Login logs user in by given user id.
	// Stay is the duration to keep the user Login state.
	// It returns a token string, use the token to keep track on the user with
	// GetUser or Logout.
	Login(id string, stay time.Duration) (string, error)
	// Logout logs the current user out. if all, logout all session of the same user.
	Logout(token string, all bool) error
	// ComparePassword cmapres the given passwod string with the hashed password of the
	// User struct.
	ComparePassword(ps string, pwd *Password) error

	// AddGroupDetail adds a group with full detail to database.
	AddGroupDetail(name string, pri []string, info *GroupInfo) (*Group, error)
	// UpdateGroupDetail updates group detail specific by id.
	UpdateGroupDetail(id string, pri []string, info *GroupInfo) error
	// FindGroup find the group specific by id.
	FindGroup(id string) (*Group, error)
	// FindGroupByName find the group specific by name.
	FindGroupByName(name string) (*Group, error)
	// FindSomeGroup find and return a slice of group specific by thier id.
	// Specific fields name for porjection select.
	FindSomeGroup(id []string, fields []string) ([]*Group, error)
	// FindAllGroup finds and return a slice of group.
	// If limit < 0 the mean using the default upper limit.
	// If limit == 0 return empty result with error indicate no result found.
	// If limit can't be greater than the default upper limit.
	// Specific fields name for porjection select.
	FindAllGroup(limit int, offsetId string, fields []string) ([]*Group, error)
	// DeleteGroup deletes a group from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	DeleteGroup(id string) error
}
