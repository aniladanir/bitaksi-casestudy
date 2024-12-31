package httphandler

import (
	"net/http"
	"strconv"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/shared/errs"
	"github.com/aniladanir/bitaksi-casestudy/shared/geojson"
	"github.com/aniladanir/bitaksi-casestudy/shared/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/shared/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type matchingHandler struct {
	logger        *zap.Logger
	driverService services.MatchingService
}

func newMatchingHandler(logger *zap.Logger, driverService services.MatchingService) *matchingHandler {
	return &matchingHandler{
		logger:        logger,
		driverService: driverService,
	}
}

func (dh *matchingHandler) FindNearestDriver(ctx fiber.Ctx) error {
	type ResponsePayload struct {
		DriverLocation *domain.DriverLocation `json:"driverLocation"`
		Distance       *domain.Distance       `json:"distance"`
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

	// call driver service
	driver, distance, err := dh.driverService.FindNearestDriverLocation(
		ctx.Context(),
		domain.UserLocation{Point: geo.(geojson.Point)},
		radius,
	)
	if err != nil {
		logger.Error("could not find nearest driver", zap.Error(err))
		if errs.IsEntityNotFoundErr(err) {
			return response.Fail(ctx, response.ErrCodeNotFound, response.ErrMsgNotFound, http.StatusNotFound)
		}
		return response.Fail(ctx, response.ErrCodeInternal, response.ErrMsgInternal, http.StatusInternalServerError)
	}

	return response.Success(ctx, &ResponsePayload{
		DriverLocation: driver,
		Distance:       distance,
	})
}
