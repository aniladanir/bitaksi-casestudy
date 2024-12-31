package httpfiber

import (
	"time"

	"github.com/aniladanir/bitaksi-casestudy/shared/log"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	HeaderXRequestID = "x-request-id"
	HeaderXTraceID   = "x-trace-id"
)

// AccessLogger middleware logs incoming http requests
func AccessLogMiddleware(logger *zap.Logger) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		log.Info("AccessLogMiddleware")
		start := time.Now()

		err := ctx.Next()

		elapsed := time.Since(start)

		CtxLogger(ctx, logger).Info("Request",
			zap.String("ip", ctx.IP()),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
			zap.Int("statusCode", ctx.Response().StatusCode()),
			zap.Duration("latency", elapsed),
			zap.Bool("success", err == nil),
		)

		return err
	}
}

// SetLogger middleware adds logger with trace abilities to the request's context
func TracingMiddleware(ctx fiber.Ctx) error {
	log.Info("TracingMiddleware")
	if ctx.Get(HeaderXTraceID) == "" {
		traceID := uuid.NewString()
		ctx.Set(HeaderXTraceID, traceID)
	}
	if ctx.Get(HeaderXRequestID) == "" {
		requestID := uuid.NewString()
		ctx.Set(HeaderXRequestID, requestID)
	}
	// store trace-id and request-id in request context
	ctx.Context().SetUserValue(CtxKeyTraceID, ctx.GetRespHeader(HeaderXTraceID))
	ctx.Context().SetUserValue(CtxKeyRequestID, ctx.GetRespHeader(HeaderXRequestID))

	return ctx.Next()
}
