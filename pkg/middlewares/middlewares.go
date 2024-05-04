package middlewares

import (
	"context"
	"log/slog"
	"os"

	"github.com/DeepAung/deep-art/pkg/prettylog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Logger() echo.MiddlewareFunc {
	slog.SetDefault(slog.New(prettylog.NewHandler(os.Stdout, nil)))

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

			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR", attrs...)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST", attrs...)
			}

			return nil
		},
	})
}
