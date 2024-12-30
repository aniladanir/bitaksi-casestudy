package geojson

const (
	TypePoint              = "Point"
	TypeLineString         = "LineString"
	TypePolygon            = "Polygon"
	TypeMultiPoint         = "MultiPoint"
	TypeMultiLineString    = "MultiLineString"
	TypeMultiPolygon       = "MultiPolygon"
	TypeGeometryCollection = "GeometryCollection"
	TypeFeature            = "Feature"
)

// Coordinate represents a single coordinate pair (longitude, latitude).
type Coordinate []float64

// Coordinates represents a set of Coordinates
type Coordinates []Coordinate

// MultiCoordinates represents a set of set of Coordinates
type MultiCoordinates []Coordinates

// Geometry is the base interface for all GeoJSON geometry objects.
type Geometry interface {
	GetType() string
	IsValid() bool
}

// Point represents a GeoJSON Point geometry.
type Point struct {
	Type        string     `json:"type"`
	Coordinates Coordinate `json:"coordinates"`
}

func (p Point) GetType() string {
	return p.Type
}

func (p Point) IsValid() bool {
	return len(p.Coordinates) == 2
}

// LineString represents a GeoJSON LineString geometry.
type LineString struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
}

func (ls LineString) GetType() string {
	return ls.Type
}

func (ls LineString) IsValid() bool {
	return len(ls.Coordinates) == 2
}

// Polygon represents a GeoJSON Polygon geometry.
type Polygon struct {
	Type        string           `json:"type"`
	Coordinates MultiCoordinates `json:"coordinates"`
}

func (p Polygon) GetType() string {
	return p.Type
}

func (p Polygon) IsValid() bool {
	if len(p.Coordinates) < 1 {
		return false
	}
	for _, ring := range p.Coordinates {
		if len(ring) < 4 {
			return false
		}
		if ring[0][0] != ring[len(ring)-1][0] || ring[0][1] != ring[len(ring)-1][1] {
			return false
		}
	}
	return true
}

// MultiPoint represents a GeoJSON MultiPoint geometry.
type MultiPoint struct {
	Type        string      `json:"type"`
	Coordinates Coordinates `json:"coordinates"`
}

func (mp MultiPoint) GetType() string {
	return mp.Type
}

func (mp MultiPoint) IsValid() bool {
	return len(mp.Coordinates) > 0
}

// MultiLineString represents a GeoJSON MultiLineString geometry.
type MultiLineString struct {
	Type        string        `json:"type"`
	Coordinates []Coordinates `json:"coordinates"`
}

func (mls MultiLineString) GetType() string {
	return mls.Type
}

func (mls MultiLineString) IsValid() bool {
	if len(mls.Coordinates) == 0 {
		return false
	}
	for _, line := range mls.Coordinates {
		if len(line) < 2 {
			return false
		}
	}
	return true
}

// MultiPolygon represents a GeoJSON MultiPolygon geometry.
type MultiPolygon struct {
	Type        string             `json:"type"`
	Coordinates []MultiCoordinates `json:"coordinates"`
}

func (mp MultiPolygon) GetType() string {
	return mp.Type
}

func (mp MultiPolygon) IsValid() bool {
	if len(mp.Coordinates) == 0 {
		return false
	}
	for _, polygon := range mp.Coordinates {
		if len(polygon) < 1 {
			return false
		}
		for _, ring := range polygon {
			if len(ring) < 4 {
				return false
			}
			if ring[0][0] != ring[len(ring)-1][0] || ring[0][1] != ring[len(ring)-1][1] {
				return false
			}
		}

	}
	return true
}

// GeometryCollection represents a GeoJSON GeometryCollection geometry.
type GeometryCollection struct {
	Type       string     `json:"type"`
	Geometries []Geometry `json:"geometries"`
}

func (gc GeometryCollection) GetType() string {
	return gc.Type
}

func (gc GeometryCollection) IsValid() bool {
	if len(gc.Geometries) == 0 {
		return false
	}
	for _, geom := range gc.Geometries {
		if !geom.IsValid() {
			return false
		}
	}
	return true
}

// Feature represents a GeoJSON Feature.
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func (f Feature) GetType() string {
	return f.Type
}

func (f Feature) IsValid() bool {
	if f.Geometry == nil {
		return false
	}
	return f.Geometry.IsValid()
}
