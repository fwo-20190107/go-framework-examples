package middleware

import (
	"io"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

var Logger *loggerMiddleware

type loggerMiddleware struct {
	logger zerolog.Logger
}

func InitLoggerMiddleware(w io.Writer) {
	Logger = &loggerMiddleware{logger: zerolog.New(w)}
}

func (m *loggerMiddleware) WithLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogMethod:    true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			c.Logger()
			m.logger.Info().
				Str("requestURI", v.URI).
				Str("method", v.Method).
				Int64("size", c.Request().ContentLength).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("remoteIP", v.RemoteIP).
				Str("agent", v.UserAgent).
				Send()
			return nil
		},
	})
}

func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}
