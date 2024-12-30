package httphandler

import (
	"net/http"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/authenticator"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type authHandler struct {
	logger        *zap.Logger
	authenticator authenticator.Authenticator
}

func newAuthHandler(logger *zap.Logger, authenticator authenticator.Authenticator) *authHandler {
	return &authHandler{
		logger:        logger,
		authenticator: authenticator,
	}
}

func (ah *authHandler) Authenticate(ctx fiber.Ctx) error {
	logger := httpfiber.CtxLogger(ctx, ah.logger)
	tokenStr := ctx.Get(fiber.HeaderAuthorization)
	if tokenStr == "" {
		logger.Error("missing authorization header")
		return response.Fail(ctx, response.ErrCodeUnauthorized, "missing authorization header", http.StatusUnauthorized)
	}
	ah.authenticator.Authenticate(ctx.Context(), tokenStr)
	authenticated, err := ah.authenticator.Authenticate(ctx.Context(), tokenStr)
	if err != nil {
		logger.Error("could not authenticate", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeInternal, response.ErrMsgInternal, http.StatusInternalServerError)
	}
	if !authenticated {
		logger.Error("token is not autheorized")
		return response.Fail(ctx, response.ErrCodeUnauthorized, "unauthorized", http.StatusUnauthorized)
	}

	return response.Success(ctx, nil)
}
