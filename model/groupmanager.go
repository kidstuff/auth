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

type BriefGroup struct {
	Id   interface{} `bson:"_id"`
	Name *string     `bson:"Name" json:",omitempty"`
}

type Group struct {
	BriefGroup `bson:"BriefGroup,inline"`
	Info       *GroupInfo `bson:"Info" json:",omitempty"`
	Privilege  []string   `bson:"Privilege" json:",omitempty"`
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
	Find(id interface{}) (*Group, error)
	// FindSome find and return a slice of group specific by thier id.
	FindSome(id ...interface{}) ([]*Group, error)
	// FindAll finds and return a slice of group.
	// If limit < 0 the mean using the default upper limit.
	// If limit == 0 return empty result with error indicate no result found.
	// If limit can't be greater than the default upper limit.
	// Specific fields name for porjection select, leave fields empty for select all.
	// offset is valid only the sort field name given and match the field type.
	FindAll(limit int, fields []string, sort string, offset interface{}) ([]*Group, error)
	// Delete deletes a group from database base on the given id;
	// It returns an error describes the first issue encountered, if any.
	Delete(id interface{}) error
}
