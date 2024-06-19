package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client = mongo.Client
type Database = mongo.Database
type Collection = mongo.Collection

type ClientOptions = options.ClientOptions
