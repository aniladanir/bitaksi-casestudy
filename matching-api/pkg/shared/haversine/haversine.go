package haversine

import (
	"fmt"
	"math"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/geojson"
)

const EarthRadiusKm = 6371.0

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// HaversineDistance calculates the distance between two GeoJSON points using the Haversine formula
func HaversineDistanceInKM(point1, point2 geojson.Point) (float64, error) {
	if point1.Type != geojson.TypePoint || point2.Type != geojson.TypePoint {
		return 0, fmt.Errorf("invalid geojson types for distance calculation")
	}

	if len(point1.Coordinates) != 2 || len(point2.Coordinates) != 2 {
		return 0, fmt.Errorf("invalid coordinates")
	}

	lat1 := DegreesToRadians(point1.Coordinates[1])
	lon1 := DegreesToRadians(point1.Coordinates[0])
	lat2 := DegreesToRadians(point2.Coordinates[1])
	lon2 := DegreesToRadians(point2.Coordinates[0])

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c, nil
}
