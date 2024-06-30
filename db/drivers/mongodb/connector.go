package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, connectionString string) (*Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
}

func SwitchToDB(ctx context.Context, client *Client, dbName string) *Database {
	return client.Database(dbName)
}

func Disconnect(ctx context.Context, client *Client) error {
	return client.Disconnect(ctx)
}

func Ping(ctx context.Context, client *Client) error {
	return client.Ping(ctx, nil)
}
