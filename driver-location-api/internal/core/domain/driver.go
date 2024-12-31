package domain

import (
	"errors"
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"github.com/google/uuid"
)

type DriverLocation struct {
	ID string `json:"id"`
	geojson.Point
}

func (dl DriverLocation) IsValid() error {
	if err := uuid.Validate(dl.ID); err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	if dl.Point.IsValid() {
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
