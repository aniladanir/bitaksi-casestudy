package httphandler

import (
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func (h *Handler) applyRoutes(accessLogger *zap.Logger) {
	// apply common middlewares
	h.app.Use(httpfiber.TracingMiddleware)
	h.app.Use(httpfiber.AccessLogMiddleware(accessLogger))

	api := h.app.Group(fmt.Sprintf("/api/%s", h.apiVersion), func(ctx fiber.Ctx) error {
		h.authHandler.Authenticate(ctx)
		return ctx.Next()
	})

	// Auth API
	authApi := api.Group("/auth")
	authApi.Post("", h.authHandler.Authenticate)

	// Driver API
	driverApi := api.Group("/driver")
	driverApi.Post("", h.driverHandler.FindNearestDriver)
}
