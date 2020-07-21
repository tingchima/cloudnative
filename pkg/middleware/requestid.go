package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// RequestIDFromContext ...
func RequestIDFromContext(ctx context.Context) string {
	value, ok := ctx.Value(echo.HeaderXRequestID).(string)
	if !ok {
		return uuid.New().String()
	}

	return value
}

// RequestID 從 ctx 中取得 request id, 如果沒有即時產生一個
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := RequestIDFromContext(c.Request().Context())

			c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			c.Response().Writer.Header().Set(echo.HeaderXRequestID, requestID)

			zlog := log.With().Str("request_id", requestID).Logger()
			ctx := zlog.WithContext(c.Request().Context())
			// ctx = errors.ContextWithXRequestID(ctx, requestID)
			ctx = context.WithValue(ctx, echo.HeaderXRequestID, requestID)
			zlog.WithContext(ctx)

			return next(c)
		}
	}
}
