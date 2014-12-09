About auth
====
[![GoDoc](http://godoc.org/github.com/kidstuff/auth?status.svg)](http://godoc.org/github.com/kidstuff/auth)  
An authentication, authorization and user management service wrtien in Go

The project provide you 2 ways to use it:  

  - A [HTTP REST API](http://kidstuff.github.io/swagger/#!/default) for authentication, authorization and user management with some client-libs for easy communication with the REST server.
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

TODO
====
* Fix bugs, improve documentation, more test (not many test right now)
* Support other type of "grant_type", become an OAuth provider.
* Better way to handle "permission" (or friendship) becom a social network.
* A complete example application (a small social network...maybe)

Community
====
https://groups.google.com/forum/#!forum/kidstuff-opensources  
We welcome any feedback, bug report, feature request or just a "Hello" from you!