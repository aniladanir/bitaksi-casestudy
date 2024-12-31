package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"github.com/aniladanir/bitaksi-casestudy/shared/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/shared/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type locationHandler struct {
	logger          *zap.Logger
	locationService services.LocationService
}

func newLocationHandler(logger *zap.Logger, locationService services.LocationService) *locationHandler {
	return &locationHandler{
		logger:          logger,
		locationService: locationService,
	}
}

func (dh *locationHandler) AddLocations(ctx fiber.Ctx) error {
	// get context logger
	logger := httpfiber.CtxLogger(ctx, dh.logger)

	// parse payload
	var payload []domain.DriverLocation
	if err := json.Unmarshal(ctx.Body(), &payload); err != nil {
		logger.Error("could not unmarshal payload", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeInvalidPayload, response.ErrMsgInvalidPayload, http.StatusBadRequest)
	}

	// validate location ids
	for i := 0; i < len(payload); i++ {
		if err := dh.locationService.IsValidID(payload[i].ID); err != nil {
			logger.Error("invalid location id", zap.Error(err), zap.Int("element", i+1))
			return response.Fail(ctx, response.ErrCodeInvalidPayload, response.ErrMsgInvalidPayload, http.StatusBadRequest)
		}
	}

	if err := dh.locationService.CreateOrUpdateDriverLocations(ctx.Context(), payload); err != nil {
		logger.Error("invalid location id", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeInternal, response.ErrMsgInternal, http.StatusInternalServerError)
	}

	return response.Success(ctx, nil)
}

func (dh *locationHandler) FindNearestDriver(ctx fiber.Ctx) error {
	type ResponseBody struct {
		Distance       domain.Distance       `json:"distance"`
		DriverLocation domain.DriverLocation `json:"location"`
	}
	// get context logger
	logger := httpfiber.CtxLogger(ctx, dh.logger)

	// parse query params
	var radius float64
	var err error
	if radius, err = strconv.ParseFloat(ctx.Query("radius"), 64); err != nil {
		logger.Error("invalid radius query param")
		return response.Fail(ctx, response.ErrCodeInvalidQueryParam, response.ErrMgInvalidQueryParam, http.StatusBadRequest)
	}

	// parse body
	geo, err := geojson.UnmarshalJSON(ctx.Body())
	if err != nil {
		logger.Error("could not unmarshal geojson data", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeInvalidPayload, response.ErrMsgInvalidPayload, http.StatusBadRequest)
	}

	// validate the type of geojson data
	if geo.GetType() != geojson.TypePoint {
		logger.Error("type of geojson data is not a point")
		return response.Fail(ctx, response.ErrCodeInvalidPayload, response.ErrMsgInvalidPayload, http.StatusBadRequest)
	}

	// validate geojson data
	if !geo.IsValid() {
		logger.Error("invalid geojson data")
		return response.Fail(ctx, response.ErrCodeInvalidPayload, response.ErrMsgInvalidPayload, http.StatusBadRequest)
	}

	// call location service
	driverLocation, distance, err := dh.locationService.FindNearestDriverDistance(
		ctx.Context(),
		domain.DriverLocation{Point: geo.(geojson.Point)},
		radius,
	)
	if err != nil {
		logger.Error("could not find driver location", zap.Error(err))
		if errs.IsInternalErr(err) {
			return response.Fail(ctx, response.ErrCodeInternal, response.ErrMsgInternal, http.StatusInternalServerError)
		}
		if errs.IsEntityNotFoundErr(err) {
			return response.Fail(ctx, response.ErrCodeNotFound, response.ErrMsgNotFound, http.StatusNotFound)
		}

	}

	return response.Success(ctx, &ResponseBody{
		Distance:       *distance,
		DriverLocation: *driverLocation,
	})
}
