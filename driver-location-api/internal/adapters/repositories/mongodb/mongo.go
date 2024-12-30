package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	return mongo.Connect(opts)
}

func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
