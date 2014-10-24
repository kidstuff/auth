About auth
====
[![GoDoc](http://godoc.org/github.com/kidstuff/auth?status.svg)](http://godoc.org/github.com/kidstuff/auth)  
An authentication, authorization and user management service wrtien in Go

The package provide you 3 ways to use it:  

  - A HTTP REST API for authentication, authorization and user management.
  - The client (currently in Go and Javascript - not finished yet [both]) for easy communication with the REST server.
  - An abstract interface (we call a manager) that let you port your Go authenthication and user management to other database. We love to support MongoDB, Google Appengine and a standar SQL but currently just a MongoDB available.
    - https://github.com/kidstuff/auth-mongo-mngr written base on mgo driver

Example
====

You can find and usage example of the package at https://github.com/kidstuff/auth-example  