package indexes

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	opts := options.CreateIndexes().SetMaxTime(1000)
	nodes := db.Collection("objects")
	_, err := nodes.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"tags", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
		{
			Keys:    bsonx.Doc{{"members.ref", bsonx.Int32(-1)}},
			Options: (options.Index()).SetBackground(true).SetSparse(true),
		},
	}, opts)
	if err != nil {
		return err
	}

	return GeoIndex(ctx, nodes, "location")
}

func GeoIndex(ctx context.Context, col *mongo.Collection, key string) error {
	_, err := col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bsonx.Doc{{
				Key: key, Value: bsonx.String("2dsphere"),
			}},
			Options: options.Index().SetSphereVersion(2).SetSparse(true).SetBackground(true),
		},
	)
	return err
}
