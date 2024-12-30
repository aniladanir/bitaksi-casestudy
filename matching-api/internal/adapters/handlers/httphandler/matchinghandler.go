package httphandler

import (
	"net/http"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/domain"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/errs"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/geojson"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/response"
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
		DriverLocation *domain.DriverLocation
		Distance       *domain.Distance
	}

	// get context logger
	logger := httpfiber.CtxLogger(ctx, dh.logger)

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
		ctx.UserContext(),
		domain.UserLocation{Point: geo.(geojson.Point)},
		ctx.Get(fiber.HeaderAuthorization),
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
