package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"github.com/aniladanir/bitaksi-casestudy/shared/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type MockLocationService struct {
	Valid          bool
	NotFound       bool
	DriverLocation domain.DriverLocation
	Distance       domain.Distance
}

func (mls *MockLocationService) FindNearestDriverDistance(ctx context.Context, location domain.DriverLocation, searchRadius float64) (*domain.DriverLocation, *domain.Distance, error) {
	if mls.NotFound {
		return nil, nil, errs.ErrEntityNotFound("not found")
	}
	return &domain.DriverLocation{
			ID: "123",
			Point: geojson.Point{
				Type: geojson.TypePoint,
				Coordinates: geojson.Coordinate{
					15.6,
					12, 5,
				},
			},
		},
		&domain.Distance{
			Distance: 140,
			Unit:     "km",
		},
		nil
}
func (*MockLocationService) CreateOrUpdateDriverLocations(ctx context.Context, locations []domain.DriverLocation) error {
	return nil
}
func (*MockLocationService) ImportLocation(ctx context.Context, reader io.Reader) error { return nil }
func (mls *MockLocationService) IsValidID(id string) error {
	if mls.Valid {
		return nil
	}
	return errors.New("error")
}

func TestFindNearestDriver(t *testing.T) {
	type ResponseBody struct {
		Distance       domain.Distance       `json:"distance"`
		DriverLocation domain.DriverLocation `json:"location"`
	}

	testCases := []struct {
		name            string
		payload         geojson.Point
		locationService services.LocationService
		expectedStatus  int
		expectedBody    response.Response
		serverPort      int
	}{
		{
			name: "should success",
			payload: geojson.Point{
				Type:        geojson.TypePoint,
				Coordinates: geojson.Coordinate{10, 10},
			},
			expectedStatus: http.StatusOK,
			locationService: &MockLocationService{
				Valid:    true,
				NotFound: false,
			},
			expectedBody: response.Response{
				Success: true,
				Data: ResponseBody{
					Distance: domain.Distance{
						Distance: 140,
						Unit:     "km",
					},
					DriverLocation: domain.DriverLocation{
						ID: "123",
						Point: geojson.Point{
							Type: geojson.TypePoint,
							Coordinates: geojson.Coordinate{
								15.6,
								12, 5,
							},
						},
					},
				},
				Code:    response.SuccessCode,
				Message: response.SuccessMsg,
			},
		},
		{
			name: "should fail due to invalid request payload",
			payload: geojson.Point{
				Type:        geojson.TypeGeometryCollection,
				Coordinates: geojson.Coordinate{10, 10},
			},
			locationService: &MockLocationService{
				Valid:    false,
				NotFound: false,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: response.Response{
				Success: false,
				Code:    response.ErrCodeInvalidPayload,
				Data:    nil,
				Message: response.ErrMsgInvalidPayload,
			},
		},
		{
			name: "should fail due to no location found",
			payload: geojson.Point{
				Type:        geojson.TypePoint,
				Coordinates: geojson.Coordinate{10, 10},
			},
			locationService: &MockLocationService{
				Valid:    true,
				NotFound: true,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: response.Response{
				Success: false,
				Code:    response.ErrCodeNotFound,
				Data:    nil,
				Message: response.ErrMsgNotFound,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locationHandler := newLocationHandler(zap.L(), tc.locationService)

			app := fiber.New()
			go func() {
				if err := app.Listen("localhost:8080"); err != nil {
					t.Errorf("could not start test server: %v", err)
				}
			}()
			defer app.Shutdown()

			app.Post("/location", locationHandler.FindNearestDriver)

			payloadBytes, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("could not marshal payload: %v", err)
			}

			resp, err := http.Post("http://localhost:8080/location?radius=1000", "application/json", bytes.NewBuffer(payloadBytes))
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status code: %d, got: %d", tc.expectedStatus, resp.StatusCode)
			}

			expectedBytes, err := json.Marshal(tc.expectedBody)
			if err != nil {
				t.Fatalf("cannot marshal expected body: %v", err)
			}
			responseBodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("cannot read response body: %v", err)
			}

			if !reflect.DeepEqual(string(expectedBytes), string(responseBodyBytes)) {
				t.Errorf("expected payload: %s, got: %s", string(expectedBytes), string(responseBodyBytes))
			}
		})
	}
}
