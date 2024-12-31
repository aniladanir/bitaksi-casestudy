package locationfinder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/shared/circuitbreaker"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"github.com/aniladanir/bitaksi-casestudy/shared/response"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type driverLocationApiClient struct {
	fasthttp.Client
	url url.URL
	cb  *circuitbreaker.CircuitBreaker
}

func NewDriverLocationApiClient(url url.URL, version string, timeout time.Duration, cb *circuitbreaker.CircuitBreaker) *driverLocationApiClient {
	return &driverLocationApiClient{
		url: *url.JoinPath(fmt.Sprintf("/api/%s", version)),
		Client: fasthttp.Client{
			WriteTimeout: timeout,
			ReadTimeout:  timeout,
		},
		cb: cb,
	}
}

func (c *driverLocationApiClient) GetNearestDriverLocation(ctx context.Context, userLocation domain.UserLocation, radius float64) (*domain.DriverLocation, *domain.Distance, error) {
	type ResponsePayload struct {
		Distance struct {
			Distance float64 `json:"distance"`
			Unit     string  `json:"unit"`
		} `json:"distance"`
		Location geojson.Point `json:"location"`
	}

	// build request url
	targetUrl := c.url.JoinPath("/driver/location")
	q := targetUrl.Query()
	q.Add("radius", strconv.FormatFloat(radius, 'f', 5, 64))
	targetUrl.RawQuery = q.Encode()

	// serialize geojson point to json
	pointJson, err := json.Marshal(userLocation)
	if err != nil {
		return nil, nil, errs.ErrInternal(err)
	}

	reqFunc := func() (any, error) {
		// build request
		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)
		req.SetRequestURI(targetUrl.String())
		req.Header.SetMethod(fasthttp.MethodPost)
		req.Header.SetContentType(fiber.MIMEApplicationJSON)
		req.Header.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
		req.SetBody(pointJson)

		// Create response object
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		// Send the request
		client := &fasthttp.Client{
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
		if err := client.Do(req, resp); err != nil {
			return nil, fmt.Errorf("could not make request: %w", err)
		}

		// handle response
		payload := &response.Response{
			Data: new(ResponsePayload),
		}
		if err := json.Unmarshal(resp.Body(), payload); err != nil {
			return false, errs.ErrInternal(fmt.Errorf("could not decode payload: %w", err))
		}
		if !payload.Success {
			switch resp.StatusCode() {
			case http.StatusNotFound:
				return nil, errs.ErrEntityNotFound("driver location")
			default:
				return nil, errs.ErrInternal(errors.New(payload.Message))
			}
		}

		data := payload.Data.(*ResponsePayload)
		return map[string]any{
			"location": &domain.DriverLocation{
				Point: data.Location,
			},
			"distance": &domain.Distance{
				Distance: data.Distance.Distance,
				Unit:     data.Distance.Unit,
			},
		}, nil
	}

	data, err := c.cb.Execute(reqFunc)
	if err != nil {
		return nil, nil, err
	}
	return data.(map[string]any)["location"].(*domain.DriverLocation), data.(map[string]any)["distance"].(*domain.Distance), nil
}
