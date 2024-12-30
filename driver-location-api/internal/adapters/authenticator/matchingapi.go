package authenticator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/errs"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Client struct {
	http.Client
	url url.URL
}

func NewMatchingApiClient(url url.URL, version string, timeout time.Duration) *Client {
	return &Client{
		url: *url.JoinPath(fmt.Sprintf("/api/%s", version)),
		Client: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Authenticate(ctx context.Context, authToken string) (bool, error) {
	// build request url
	targetUrl := c.url.JoinPath("/auth")

	// build request
	req := &http.Request{
		Method: http.MethodGet,
		URL:    targetUrl,
		Header: http.Header{
			fiber.HeaderAuthorization:  []string{authToken},
			httpfiber.HeaderXTraceID:   []string{ctx.Value(httpfiber.CtxKeyTraceID).(string)},
			httpfiber.HeaderXRequestID: []string{uuid.NewString()},
		},
	}
	req = req.WithContext(ctx)

	// do request
	resp, err := c.Do(req)
	if err != nil {
		return false, errs.ErrInternal(fmt.Errorf("could not make request: %w", err))
	}

	// handle response
	payload := new(response.Response)
	if err := json.NewDecoder(resp.Body).Decode(payload); err != nil {
		return false, errs.ErrInternal(fmt.Errorf("could not decode payload: %w", err))
	}
	if !payload.Success {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return false, nil
		default:
			return false, errs.ErrInternal(errors.New(payload.Message))
		}
	}

	return true, nil
}
