About auth
====
[![GoDoc](http://godoc.org/github.com/kidstuff/auth?status.svg)](http://godoc.org/github.com/kidstuff/auth)  
An authentication, authorization and user management service wrtien in Go

The package provide you 3 ways to use it:  

  - A [HTTP REST API](http://kidstuff.github.io/swagger/#!/default) for authentication, authorization and user management.
  - The clients for easy communication with the REST server.
    - https://github.com/kidstuff/auth-angular-client AngualrJS client.
  - An abstract interface (we call a manager) that let you port your Go authenthication and user management to other database. We love to support MongoDB, Google Appengine and a standar SQL but currently just a MongoDB available.
    - https://github.com/kidstuff/auth-mongo-mngr written base on mgo driver

Documentation
====
Please prefer:
* https://github.com/kidstuff/auth/wiki
* http://godoc.org/github.com/kidstuff/auth
* http://kidstuff.github.io/swagger/#!/default
Each sub project may have their own document.

Example
====
You can find and usage example of the package at https://github.com/kidstuff/auth-example  

Community
====
https://groups.google.com/forum/#!forum/kidstuff-opensources