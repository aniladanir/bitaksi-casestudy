package services

import (
	"context"
	"fmt"
	"io"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/importer"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/repositories"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"github.com/aniladanir/bitaksi-casestudy/shared/haversine"
)

type LocationService interface {
	CreateOrUpdateDriverLocations(ctx context.Context, locations []domain.DriverLocation) error
	FindNearestDriverDistance(ctx context.Context, location domain.DriverLocation, searchRadius float64) (*domain.DriverLocation, *domain.Distance, error)
	ImportLocation(ctx context.Context, reader io.Reader) error
	IsValidID(id string) error
}

type locationService struct {
	locationImporter importer.Importer
	locationRepo     repositories.LocationRepository
}

func NewLocationService(repo repositories.LocationRepository, locationImporter importer.Importer) *locationService {
	return &locationService{
		locationImporter: locationImporter,
		locationRepo:     repo,
	}
}

func (ls *locationService) CreateOrUpdateDriverLocations(ctx context.Context, locations []domain.DriverLocation) error {
	return ls.locationRepo.UpsertMany(ctx, locations)
}

func (ls *locationService) IsValidID(id string) error {
	return ls.locationRepo.IsValidID(id)
}

func (ls *locationService) FindNearestDriverDistance(ctx context.Context, userLocation domain.DriverLocation, searchRadius float64) (*domain.DriverLocation, *domain.Distance, error) {
	driverLocation, err := ls.locationRepo.GetNearestDriverLocation(ctx, userLocation, searchRadius)
	if err != nil {
		return nil, nil, err
	}

	distanceKM, err := haversine.HaversineDistanceInKM(driverLocation.Point, userLocation.Point)
	if err != nil {
		return nil, nil, errs.ErrInternal(fmt.Errorf("could not calculate distance between points: %w", err))
	}

	return driverLocation, &domain.Distance{
		Distance: distanceKM,
		Unit:     "km",
	}, nil
}

func (ls *locationService) ImportLocation(ctx context.Context, reader io.Reader) error {
	return ls.locationImporter.ImportCoordinates(ctx, reader)
}
