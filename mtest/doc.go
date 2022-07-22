/*
Package mtest is a helper library for integration test with some MongoDB instance.

This package provide a clean and readable way to ensure testing pre-condition and
validate post-condition for some test case.

As the design of MongoDB official Golang library, the document to be checked by this
package is given by bson.M, bson.D or any Golang struct (with "bson" structure tag).

See https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo for more details.
*/
package mtest
