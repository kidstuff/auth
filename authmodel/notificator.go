// Copyright 2012 The KidStuff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authmodel

var DEFAULT_NOTIFICATOR Notificator

type Notificator interface {
	SendMail(subject, message, from, to string) error
}
