package indexes

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func CreateIndexes(ctx context.Context, db *mongo.Database) error {
	opts := mongoOptions.CreateIndexes().SetMaxTime(1000)
	objects := db.Collection("objects")
	_, err := objects.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"tags", bsonx.Int32(-1)}},
			Options: (mongoOptions.Index()).SetSparse(true),
		},
		{
			Keys:    bsonx.Doc{{"members.ref", bsonx.Int32(-1)}},
			Options: (mongoOptions.Index()).SetSparse(true),
		},
	}, opts)
	if err != nil {
		return err
	}

	return GeoIndex(ctx, objects, "location")
}

func GeoIndex(ctx context.Context, col *mongo.Collection, key string) error {
	_, err := col.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bsonx.Doc{{
				Key: key, Value: bsonx.String("2dsphere"),
			}},
			Options: mongoOptions.Index().SetSphereVersion(2).SetSparse(true),
		},
	)
	return err
}
