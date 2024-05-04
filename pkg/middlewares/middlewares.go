package middlewares

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

// RemoteIP string
// Host string
func (m *Middleware) Logger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("uri", v.URI),
				slog.String("method", v.Method),
				slog.String("latency", v.Latency.String()),
				slog.Int("status", v.Status),
			}

			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", attrs...)
			} else {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR", attrs...)
			}
			return nil
		},
	})
}
