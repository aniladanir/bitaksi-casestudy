package mongodb

import (
	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverLocation struct {
	ID       primitive.ObjectID `bson:"_id"`
	Location geojson.Point      `bson:"location"`
}
