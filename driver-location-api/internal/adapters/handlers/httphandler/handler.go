package httphandler

import (
	"time"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/authenticator"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/services"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Handler struct {
	app             *fiber.App
	apiVersion      string
	logger          *zap.Logger
	locationHandler *locationHandler
	authHandler     *authHandler
}

type ServerConfig struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

func NewHandler(serverCfg ServerConfig, logger *zap.Logger, accessLogger *zap.Logger, locationService services.LocationService, authenticator authenticator.Authenticator, apiVersion string) *Handler {
	h := &Handler{
		app: fiber.New(fiber.Config{
			ReadTimeout:  serverCfg.ReadTimeout,
			WriteTimeout: serverCfg.WriteTimeout,
			IdleTimeout:  serverCfg.IdleTimeout,
		}),
		logger:          logger,
		locationHandler: newLocationHandler(logger.With(zap.String("handler", "location")), locationService),
		authHandler:     newAuthHandler(logger.With(zap.String("handler", "auth")), authenticator),
	}
	h.applyRoutes(accessLogger)
	return h
}

func (h *Handler) Listen(address string) error {
	err := h.app.Listen(address)
	h.logger.Error("listen on address failed", zap.Error(err), zap.String("address", address))
	return err
}

func (h *Handler) Shutdown() error {
	if err := h.app.Shutdown(); err != nil {
		h.logger.Error("shutdown failed", zap.Error(err))
		return err
	}
	return nil
}
