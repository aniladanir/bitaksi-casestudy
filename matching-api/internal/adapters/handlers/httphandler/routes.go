package httphandler

import (
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/shared/httpfiber"
	"go.uber.org/zap"
)

func (h *Handler) applyRoutes(accessLogger *zap.Logger) {
	// apply common middlewares
	h.app.Use(httpfiber.TracingMiddleware)
	h.app.Use(httpfiber.AccessLogMiddleware(accessLogger))

	api := h.app.Group(fmt.Sprintf("/api/%s", h.apiVersion))

	// Auth API
	api.Post("/auth", h.authHandler.Authenticate(false))

	// Driver API
	driverApi := api.Group("/match", h.authHandler.Authenticate(true))
	driverApi.Post("/driver", h.driverHandler.FindNearestDriver)
}
