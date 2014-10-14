// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"errors"
)

var (
	ErrDuplicateName = errors.New("auth: duplicate Group Name")
)

type Group struct {
	Id        *string    `bson:"-"`
	Name      *string    `bson:"Name,omitempty" json:",omitempty"`
	Info      *GroupInfo `bson:"Info,omitempty" json:",omitempty"`
	Privilege []string   `bson:"Privilege,omitempty" json:",omitempty"`
}

type GroupInfo struct {
	Description *string `bson:"Description" json:",omitempty"`
}

type GroupManager interface {
	// AddDetail adds a group with full detail to database.
	AddDetail(*Group) (*Group, error)
	// UpdateDetail updates group detail specific by id.
	UpdateDetail(*Group) error
	// Find find the group specific by id.
	Find(id string) (*Group, error)
	// FindByName find the group specific by name.
	FindByName(name string) (*Group, error)
	// FindSome find and return a slice of group specific by thier id.
	FindSome(id ...string) ([]*Group, error)
	// FindAll finds and return a slice of group.
	// If limit < 0 the mean using the default upper limit.
	// If limit == 0 return empty result with error indicate no result found.
	// If limit can't be greater than the default upper limit.
	// Specific fields name for porjection select.
	FindAll(limit int, offsetId string, fields []string) ([]*Group, error)
	// Delete deletes a group from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	Delete(id string) error
}
