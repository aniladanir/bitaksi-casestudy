package locationfinder

import (
	"context"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
)

type LocationFinder interface {
	GetNearestDriverLocation(ctx context.Context, userLocation domain.UserLocation, radius float64) (*domain.DriverLocation, *domain.Distance, error)
}
