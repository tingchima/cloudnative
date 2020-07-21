package middleware

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Recovery handling panic error
func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					trace := make([]byte, 4096)
					runtime.Stack(trace, true)
					customFields := map[string]interface{}{
						"uri":         c.Request().RequestURI,
						"stack_error": string(trace),
					}
					err, ok := r.(error)
					if !ok {
						if err == nil {
							err = fmt.Errorf("%v", r)
						} else {
							err = fmt.Errorf("%v", err)
						}
					}
					logger := log.With().Fields(customFields).Logger()
					logger.Error().Msgf("http: unknown error: %v", err)

					// _ = c.JSON(500, errors.NewWithMessage(errors.ErrInternalError, err.Error()))
					c.JSON(500, gin.H{})
				}
			}()
			return next(c)
		}
	}
}
