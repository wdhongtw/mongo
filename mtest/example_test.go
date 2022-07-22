package mtest

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Example_basic() {
	t := &testing.T{}
	someCollection := &mongo.Collection{}

	// Prepare pre-condition with Condition.
	For(t, someCollection).Condition().
		Empty().
		Exists(bson.M{"name": "Alice", "age": 22})

	// Some manipulation on the document ...

	// Check post-condition with Assert.
	For(t, someCollection).Assert().
		Exists(bson.M{"name": "Alice", "age": 24}).
		NotExists(bson.M{"name": "Bob"})

	// Can also use Golang struct for checking clause
	type Person struct {
		Name string `bson:"name"`
		Age  int    `bson:"age"`
	}
	For(t, someCollection).Assert().
		Exists(Person{Name: "Alice", Age: 22})

	// Use "From" helper to get collection if we only have the client object
	client, _ := mongo.NewClient()
	For(t, From(client, "company", "person")).Condition().
		Empty().
		Exists(bson.M{"name": "Alice", "age": 22})

	// Use "OfId" helper to get a filter for some document ID
	For(t, From(client, "company", "person")).Condition().
		Empty().
		NotExists(OfId("some-id"))
}

func Example_noExistence() {
	t := &testing.T{}
	client, _ := mongo.NewClient()

	// Cleanup the collection and insert one document.
	For(t, From(client, "company", "person")).Condition().
		Empty().
		Exists(bson.M{"name": "Alice", "age": 22})

	// Delete the document ...

	// Check post-condition with Assert.
	For(t, From(client, "company", "person")).Assert().
		NotExists(bson.M{"name": "Alice", "age": 24})

	// If we need to check the collection is indeed cleaned up.
	For(t, From(client, "company", "person")).Assert().
		Empty()
}

func Example_multipleRequirement() {
	t := &testing.T{}
	client, _ := mongo.NewClient()

	// Batch insert logic ...

	// Can check multiple document at once.
	For(t, From(client, "company", "person")).Assert().
		Exists(
			bson.M{"name": "Alice", "age": 24},
			bson.M{"name": "Bob", "age": 16},
			bson.M{"name": "Charlie", "age": 53},
		)

	// Another checking style.
	For(t, From(client, "company", "person")).Assert().
		Exists(bson.M{"name": "Alice", "age": 24}).
		Exists(bson.M{"name": "Bob", "age": 16}).
		Exists(bson.M{"name": "Charlie", "age": 53})
}

func Example_immediatelyFail() {
	t := &testing.T{}
	client, _ := mongo.NewClient()

	// First part of logic being testing ...

	// Check the post-condition and return immediately
	// if the condition doesn't meet the requirement
	For(t, From(client, "company", "person")).Require().
		NotExists(bson.M{"name": "Alice", "age": 24})

	// Second part of logic being testing ...

	For(t, From(client, "company", "person")).Require().
		Exists(bson.M{"name": "Alice", "age": 24})
}
