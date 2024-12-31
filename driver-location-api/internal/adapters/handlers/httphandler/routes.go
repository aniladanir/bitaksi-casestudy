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

	// Driver API
	driverApi := api.Group("/driver")
	driverApi.Put("/location", h.locationHandler.AddLocations)
	driverApi.Post("/location", h.locationHandler.FindNearestDriver)
}
