package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/repositories/mongodb"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LocationRepository interface {
	UpsertMany(ctx context.Context, locations []domain.DriverLocation) error
	GetNearestDriverLocation(ctx context.Context, userLocation domain.DriverLocation, radius float64) (*domain.DriverLocation, error)
	IsValidID(id string) error
}

type locationRepository struct {
	driverLocationDB *mongo.Collection
}

func NewLocationRepository(mongoDB *mongo.Database) *locationRepository {
	return &locationRepository{
		driverLocationDB: mongoDB.Collection("driver-location"),
	}
}

func (lr *locationRepository) IsValidID(id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errs.ErrInternal(fmt.Errorf("invalid object id: %w", err))
	}
	return nil
}

func (lr *locationRepository) UpsertMany(ctx context.Context, locations []domain.DriverLocation) error {
	var err error
	models := make([]mongo.WriteModel, 0, len(locations))
	for _, l := range locations {
		var objectID primitive.ObjectID
		if l.ID == "" {
			objectID = primitive.NewObjectID()
		} else {
			objectID, err = primitive.ObjectIDFromHex(l.ID)
			if err != nil {
				return errs.ErrInternal(fmt.Errorf("invalid location id: %w", err))
			}
		}
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": objectID}).
			SetUpdate(bson.M{"$set": bson.M{
				"_id":      objectID,
				"location": l.Point,
			}}).
			SetUpsert(true)
		models = append(models, model)
	}

	if _, err := lr.driverLocationDB.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false)); err != nil {
		return errs.ErrInternal(fmt.Errorf("failed to bulk write: %w", err))
	}

	return nil
}

func (lr *locationRepository) GetNearestDriverLocation(ctx context.Context, location domain.DriverLocation, radius float64) (*domain.DriverLocation, error) {
	filter := bson.M{"location": bson.M{
		"$near": bson.M{
			"$geometry": bson.M{
				"type":        location.Type,
				"coordinates": location.Coordinates,
			},
			"$maxDistance": radius,
		},
	},
	}

	var result mongodb.DriverLocation
	if err := lr.driverLocationDB.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.ErrEntityNotFound("driver location")
		}
		return nil, errs.ErrInternal(err)
	}

	return &domain.DriverLocation{
		Point: result.Location,
	}, nil
}
