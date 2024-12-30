package domain

import (
	"errors"
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/geojson"
	"github.com/google/uuid"
)

type DriverLocation struct {
	ID       string `json:"id"`
	Location geojson.Point
}

func (dl DriverLocation) IsValid() error {
	if err := uuid.Validate(dl.ID); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if dl.Location.IsValid() {
		return errors.New("invalid geojson point")
	}
	return nil
}

type Distance struct {
	Distance float64 `json:"distance"`
	Unit     string  `json:"unit"`
}

func (dl Distance) IsValid() error {
	if dl.Distance < 0 {
		return errors.New("distance cannot be negative")
	}
	return nil
}
