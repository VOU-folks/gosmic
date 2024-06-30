package mongodb

import (
	"context"
)

func GetCollection(ctx context.Context, db *Database, name string) *Collection {
	return db.Collection(name)
}

func InsertOne(ctx context.Context, doc interface{}, collection *Collection) (*InsertOneResult, error) {
	return collection.InsertOne(ctx, doc)
}

func InsertMany(ctx context.Context, docs []interface{}, collection *Collection) (*InsertManyResult, error) {
	return collection.InsertMany(ctx, docs)
}
