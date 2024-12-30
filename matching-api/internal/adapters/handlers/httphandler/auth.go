package httphandler

import (
	"net/http"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type authHandler struct {
	logger *zap.Logger
}

func newAuthHandler(logger *zap.Logger) *authHandler {
	return &authHandler{
		logger: logger,
	}
}

func (ah *authHandler) Authenticate(ctx fiber.Ctx) error {
	logger := httpfiber.CtxLogger(ctx, ah.logger)
	tokenStr := ctx.Get(fiber.HeaderAuthorization)
	if tokenStr == "" {
		logger.Error("missing authorization header")
		return response.Fail(ctx, response.ErrCodeUnauthorized, "missing authorization header", http.StatusUnauthorized)
	}
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &JwtClaims{})
	if err != nil {
		logger.Error("token is in unknown format", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeUnauthorized, "token is in unknown format", http.StatusUnauthorized)
	}

	if err = token.Claims.Valid(); err != nil {
		logger.Error("token validation failed", zap.Error(err))
		return response.Fail(ctx, response.ErrCodeUnauthorized, err.Error(), http.StatusUnauthorized)
	}

	return response.Success(ctx, nil)
}
