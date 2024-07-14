package indexes

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	err := CreateIndexesForNodes(ctx, db)
	if err != nil {
		return err
	}

	err = CreateIndexesForWays(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

func CreateIndexesForNodes(ctx context.Context, db *mongo.Database) error {
	nodes := db.Collection("nodes")
	_, err := nodes.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{"tags", -1}},
			Options: (mongoOptions.Index()).SetSparse(true),
		},
		mongoOptions.CreateIndexes().SetMaxTime(1000),
	)
	if err != nil {
		return err
	}

	return GeoIndex(ctx, nodes)
}

func CreateIndexesForWays(ctx context.Context, db *mongo.Database) error {
	ways := db.Collection("ways")
	_, err := ways.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{"tags", -1}},
			Options: (mongoOptions.Index()).SetSparse(true),
		},
		mongoOptions.CreateIndexes().SetMaxTime(1000),
	)
	if err != nil {
		return err
	}

	return GeoIndex(ctx, ways)
}

func GeoIndex(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{"location", "2dsphere"}},
			Options: mongoOptions.Index().SetSphereVersion(2).SetSparse(true),
		},
	)
	return err
}
