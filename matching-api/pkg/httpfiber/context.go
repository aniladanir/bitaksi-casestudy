package httpfiber

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	CtxKeyTraceID   = "trace-id"
	CtxKeyRequestID = "request-id"
)

func CtxLogger(ctx fiber.Ctx, logger *zap.Logger) *zap.Logger {
	traceID := ctx.Get(HeaderXTraceID)
	requestID := ctx.Get(HeaderXRequestID)
	return logger.With(zap.String(CtxKeyTraceID, traceID)).With(zap.String(CtxKeyRequestID, requestID))
}
