package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	return mongo.Connect(opts)
}

func CreateDatabase(ctx context.Context, client *mongo.Client, dbName string) (*mongo.Database, error) {
	db := client.Database(dbName)
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"location": "2dsphere"},
		Options: options.Index().SetName("location_2dsphere"),
	}
	if _, err := db.Collection("driver-location").Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, fmt.Errorf("could not create 2dsphere index on driver-location collection: %w", err)
	}
	return client.Database(dbName), nil
}
