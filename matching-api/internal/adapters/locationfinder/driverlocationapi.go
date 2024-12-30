package locationfinder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/errs"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/geojson"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type driverLocationApiClient struct {
	http.Client
	url url.URL
}

func NewDriverLocationApiClient(url url.URL, version string, timeout time.Duration) *driverLocationApiClient {
	return &driverLocationApiClient{
		url: *url.JoinPath(fmt.Sprintf("/api/%s", version)),
		Client: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *driverLocationApiClient) GetNearestDriverLocation(ctx context.Context, userLocation domain.UserLocation, authToken string) (*domain.DriverLocation, *domain.Distance, error) {
	type ResponsePayload struct {
		Distance struct {
			Distance float64 `json:"distance"`
			Unit     string  `json:"unit"`
		} `json:"distance"`
		Location geojson.Point `json:"location"`
	}

	// build request url
	targetUrl := c.url.JoinPath("/driver/location")
	q := c.url.Query()
	q.Add("type", "nearest")
	targetUrl.RawQuery = q.Encode()

	// serialize geojson point to json
	pointJson, err := json.Marshal(userLocation)
	if err != nil {
		return nil, nil, errs.ErrInternal(err)
	}

	// build request
	req := &http.Request{
		Method: http.MethodPost,
		URL:    targetUrl,
		Header: http.Header{
			fiber.HeaderAuthorization:  []string{authToken},
			httpfiber.HeaderXTraceID:   []string{ctx.Value(httpfiber.CtxKeyTraceID).(string)},
			httpfiber.HeaderXRequestID: []string{uuid.NewString()},
		},
		Body: io.NopCloser(bytes.NewReader(pointJson)),
	}
	req = req.WithContext(ctx)

	// do request
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("could not make request: %w", err)
	}

	// handle response
	payload := &response.Response{
		Data: new(ResponsePayload),
	}
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return nil, nil, fmt.Errorf("could not decode payload: %w", err)
	}
	if !payload.Success {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, nil, errs.ErrEntityNotFound("driver location")
		default:
			return nil, nil, errs.ErrInternal(errors.New(payload.Message))
		}
	}

	data := payload.Data.(*ResponsePayload)

	return &domain.DriverLocation{
			Point: data.Location,
		}, &domain.Distance{
			Distance: data.Distance.Distance,
			Unit:     data.Distance.Unit,
		}, nil
}
