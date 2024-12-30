package importer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/repositories"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/errs"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/geojson"
)

const (
	KeyLatitude   = "latitude"
	KeyLongtitude = "longtitude"
)

type csvImporter struct {
	locationRepo repositories.LocationRepository
}

func NewCsvImporter(locationRepo repositories.LocationRepository) *csvImporter {
	return &csvImporter{
		locationRepo: locationRepo,
	}
}

func (ci *csvImporter) ImportCoordinates(ctx context.Context, reader io.Reader) error {
	csvReader := csv.NewReader(reader)

	// reader header
	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("encountered error when reading header: %w", err)
	}
	longtitudeIndex, latitudeIndex, err := ci.validateHeader(header)
	if err != nil {
		return err
	}

	locations := make([]domain.DriverLocation, 0, 512)
	rowNumber := 1
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errs.ErrInternal(fmt.Errorf("encountered error when reading %d row: %w", rowNumber, err))
		}
		if len(row) != 2 {
			return errs.ErrInternal(fmt.Errorf("expected number of columns on row %d is incorrect", rowNumber))
		}

		longtitude, err := strconv.ParseFloat(row[longtitudeIndex], 64)
		if err != nil {
			return errs.ErrInternal(fmt.Errorf("could not parse longtitude on row %d: %w", rowNumber, err))
		}
		latitude, err := strconv.ParseFloat(row[latitudeIndex], 64)
		if err != nil {
			return errs.ErrInternal(fmt.Errorf("could not parse latitude on row %d: %w", rowNumber, err))
		}

		location := domain.DriverLocation{
			Location: geojson.Point{
				Type:        geojson.TypePoint,
				Coordinates: geojson.Coordinate{longtitude, latitude},
			},
		}

		locations = append(locations, location)

		rowNumber++
	}

	if len(locations) > 0 {
		if err := ci.locationRepo.UpsertMany(ctx, locations); err != nil {
			return errs.ErrInternal(err)
		}
	}

	return nil
}

func (ci *csvImporter) validateHeader(header []string) (longtitudeIndex, latitudeIndex int, err error) {
	if len(header) != 2 {
		err = errors.New("expecting number of columns is two")
		return
	}
	for i, h := range header {
		if strings.ToLower(h) != KeyLatitude || strings.ToLower(h) != KeyLongtitude {
			err = fmt.Errorf("unknown column name: %s", h)
			return
		}
		if strings.ToLower(h) == KeyLongtitude {
			longtitudeIndex = i
		}
		if strings.ToLower(h) == KeyLatitude {
			latitudeIndex = i
		}
	}
	return
}
