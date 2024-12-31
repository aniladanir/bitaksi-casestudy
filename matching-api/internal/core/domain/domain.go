package domain

import (
	"errors"

	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
)

// Distance Unit Types
const (
	UnitDistanceKiloMeter = "km"
	UnitDistanceMeter     = "m"
)

type DriverLocation struct {
	geojson.Point
}

func (dl DriverLocation) IsValid() error {
	if !dl.Point.IsValid() {
		return errors.New("invalid geojson point data")
	}
	return nil
}

type UserLocation struct {
	geojson.Point
}

func (ul UserLocation) IsValid() error {
	if !ul.Point.IsValid() {
		return errors.New("invalid geojson point data")
	}
	return nil
}

type Distance struct {
	Distance float64 `json:"distance"`
	Unit     string  `json:"unit"`
}

func (d Distance) IsValid() error {
	if d.Distance < 0 {
		return errors.New("negative distance")
	}
	switch d.Unit {
	case UnitDistanceKiloMeter, UnitDistanceMeter:
		return nil
	default:
		return errors.New("unknown distance unit")
	}
}
