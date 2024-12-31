package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

// MarshalJSON serializes a Geometry object into a GeoJSON string.
func MarshalJSON(g Geometry) (string, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return "", fmt.Errorf("error marshalling json: %w", err)
	}
	return string(b), nil
}

// UnmarshalJSON deserializes a GeoJSON string to respective struct
func UnmarshalJSON(data []byte) (Geometry, error) {
	var geoType struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(data, &geoType)
	if err != nil {
		return nil, err
	}

	switch geoType.Type {
	case TypePoint:
		var p Point
		err := json.Unmarshal(data, &p)
		if err != nil {
			return nil, err
		}
		return p, nil
	case TypeLineString:
		var ls LineString
		err := json.Unmarshal(data, &ls)
		if err != nil {
			return nil, err
		}
		return ls, nil
	case TypePolygon:
		var polygon Polygon
		err := json.Unmarshal(data, &polygon)
		if err != nil {
			return nil, err
		}
		return polygon, nil
	case TypeMultiPoint:
		var mp MultiPoint
		err := json.Unmarshal(data, &mp)
		if err != nil {
			return nil, err
		}
		return mp, nil
	case TypeMultiLineString:
		var mls MultiLineString
		err := json.Unmarshal(data, &mls)
		if err != nil {
			return nil, err
		}
		return mls, nil
	case TypeMultiPolygon:
		var mp MultiPolygon
		err := json.Unmarshal(data, &mp)
		if err != nil {
			return nil, err
		}
		return mp, nil
	case TypeGeometryCollection:
		var gc GeometryCollection
		err := json.Unmarshal(data, &gc)
		if err != nil {
			return nil, err
		}
		return gc, nil
	case TypeFeature:
		var f Feature
		err := json.Unmarshal(data, &f)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("invalid geometry type")
	}
}
