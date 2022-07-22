package mtest

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// For binds a Mongo collection with a test reporting interface.
func For(t *testing.T, col *mongo.Collection) *Collection {
	return &Collection{
		t:   t,
		col: col,
	}
}

// From is a shortcut for construct a collection from client object.
func From(client *mongo.Client, database, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}

// OfId return a document filter which matches a document by document ID ("_id" field).
func OfId(id interface{}) interface{} {
	return bson.M{"_id": id}
}

type Collection struct {
	col *mongo.Collection
	t   *testing.T
}

func (c *Collection) Condition() *Condition {
	return &Condition{
		t:   c.t,
		col: c.col,
	}
}

// Condition is a object for ensure(t.FailNow) pre-condition for some test case.
type Condition struct {
	col *mongo.Collection
	t   *testing.T
}

func (c *Collection) Require() *Require {
	return &Require{
		t:   c.t,
		col: c.col,
	}
}

// Require is a object for validate(t.FailNow) post-condition for some test case.
type Require struct {
	col *mongo.Collection
	t   *testing.T
}

func (c *Collection) Assert() *Assert {
	return &Assert{
		t:   c.t,
		col: c.col,
	}
}

// Assert is a object for validate(t.Fail) post-condition for some test case.
type Assert struct {
	col *mongo.Collection
	t   *testing.T
}

func (c *Condition) Exists(documents ...interface{}) *Condition {
	conditionExists(c.t, c.col, documents...)
	return c
}

func (c *Condition) NotExists(documents ...interface{}) *Condition {
	conditionNotExists(c.t, c.col, documents...)
	return c
}

func (c *Condition) Empty() *Condition {
	conditionEmpty(c.t, c.col)
	return c
}

func (r *Require) Exists(documents ...interface{}) *Require {
	requireExists(r.t, r.col, documents...)
	return r
}

func (r *Require) NotExists(documents ...interface{}) *Require {
	requireNotExists(r.t, r.col, documents...)
	return r
}

func (r *Require) Empty() *Require {
	requireEmpty(r.t, r.col)
	return r
}

func (r *Assert) Exists(documents ...interface{}) *Assert {
	assertExists(r.t, r.col, documents...)
	return r
}

func (r *Assert) NotExists(documents ...interface{}) *Assert {
	assertNotExists(r.t, r.col, documents...)
	return r
}

func (r *Assert) Empty() *Assert {
	assertEmpty(r.t, r.col)
	return r
}

func conditionExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		_, err := col.UpdateOne(context.TODO(), document, options.Update().SetUpsert(true))
		if err != nil {
			t.Fatalf("can not ensure document [%#v] in collection [%v]: %v", document, col.Name(), err)
		}
	}
}

func conditionNotExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		result := col.FindOne(context.TODO(), document)
		if result.Err() == nil {
			t.Fatalf("document [%#v] found in collection [%v]", document, col.Name())
		}
	}
}

func conditionEmpty(t *testing.T, col *mongo.Collection) {
	err := col.Drop(context.TODO())
	if err != nil {
		t.Fatalf("can not drop the collection [%v]: %v", col.Name(), err)
	}
}

func requireExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		result := col.FindOne(context.TODO(), document)
		if result.Err() != nil {
			t.Fatalf("can not check existence for document [%#v] in collection [%v]: %v", document, col.Name(), result.Err())
		}
	}
}

func requireNotExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		result := col.FindOne(context.TODO(), document)
		if result.Err() == nil {
			t.Fatalf("document [%#v] found in collection [%v]", document, col.Name())
		}
	}
}

func requireEmpty(t *testing.T, col *mongo.Collection) {
	count, err := col.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("can not count the collection [%v]: %v", col.Name(), err)
	}
	if count != 0 {
		t.Fatalf("still has [%v] document in collection [%v]", count, col.Name())
	}
}

func assertExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		result := col.FindOne(context.TODO(), document)
		if result.Err() != nil {
			t.Errorf("can not check existence for document [%#v] in collection [%v]: %v", document, col.Name(), result.Err())
		}
	}
}

func assertNotExists(t *testing.T, col *mongo.Collection, documents ...interface{}) {
	for _, document := range documents {
		result := col.FindOne(context.TODO(), document)
		if result.Err() == nil {
			t.Errorf("document [%#v] found in collection [%v]", document, col.Name())
		}
	}
}

func assertEmpty(t *testing.T, col *mongo.Collection) {
	count, err := col.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		t.Fatalf("can not count the collection [%v]: %v", col.Name(), err)
	}
	if count != 0 {
		t.Errorf("still has [%v] document in collection [%v]", count, col.Name())
	}
}
