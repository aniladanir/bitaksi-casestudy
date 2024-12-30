package httphandler

import (
	"fmt"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/httpfiber"
	"go.uber.org/zap"
)

func (h *Handler) applyRoutes(accessLogger *zap.Logger) {
	// apply common middlewares
	h.app.Use(httpfiber.TracingMiddleware)
	h.app.Use(httpfiber.AccessLogMiddleware(accessLogger))

	api := h.app.Group(fmt.Sprintf("/api/%s", h.apiVersion), h.authHandler.Authenticate)

	// Distance API
	distanceApi := api.Group("/distance")
	distanceApi.Post("/point", h.locationHandler.FindDriverDistanceToPoint)
}
