package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/prettylog"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Middleware struct {
	usersSvc *services.UsersSvc
	cfg      *config.Config
}

func NewMiddleware(usersSvc *services.UsersSvc, cfg *config.Config) *Middleware {
	return &Middleware{
		usersSvc: usersSvc,
		cfg:      cfg,
	}
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
				slog.String("method", c.Request().Method),
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

func (m *Middleware) OnlyAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("accessToken")
		if err != nil {
			utils.ClearCookies(c)
			c.Redirect(http.StatusMovedPermanently, "/signin")
			return nil
		}

		tokenString := cookie.Value
		if tokenString == "" {
			utils.ClearCookies(c)
			c.Redirect(http.StatusMovedPermanently, "/signin")
			return nil
		}

		claims, err := mytoken.ParseToken(mytoken.Access, m.cfg.Jwt.SecretKey, tokenString)
		if err != nil {
			utils.ClearCookies(c)
			c.Redirect(http.StatusMovedPermanently, "/signin")
			return nil
		}

		has, err := m.usersSvc.HasAccessToken(claims.Payload.UserId, tokenString)
		if err != nil || !has {
			utils.ClearCookies(c)
			c.Redirect(http.StatusMovedPermanently, "/signin")
			return nil
		}

		c.Set("payload", claims.Payload)

		return next(c)
	}
}
