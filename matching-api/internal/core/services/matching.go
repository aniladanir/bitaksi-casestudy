package services

import (
	"context"
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/adapters/locationfinder"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
)

type MatchingService interface {
	FindNearestDriverLocation(ctx context.Context, userLocation domain.UserLocation, radius float64) (*domain.DriverLocation, *domain.Distance, error)
}

type matchingService struct {
	locationFinder locationfinder.LocationFinder
}

func NewDriverService(locationFinder locationfinder.LocationFinder) *matchingService {
	return &matchingService{
		locationFinder: locationFinder,
	}
}

func (ds *matchingService) FindNearestDriverLocation(ctx context.Context, userLocation domain.UserLocation, radius float64) (*domain.DriverLocation, *domain.Distance, error) {
	driverLocation, distanceToUser, err := ds.locationFinder.GetNearestDriverLocation(ctx, userLocation, radius)
	if err != nil {
		return nil, nil, err
	}
	if err := driverLocation.IsValid(); err != nil {
		return nil, nil, errs.ErrInternal(fmt.Errorf("got invalid data from location finder: %w", err))
	}
	return driverLocation, distanceToUser, nil
}
