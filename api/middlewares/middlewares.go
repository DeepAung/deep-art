package middlewares

import (
	"context"
	"errors"
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

type Authorized struct {
	c       echo.Context
	mid     *Middleware
	payload mytoken.Payload
}

type AuthorizedOpt func(*Authorized) error

func SetPayload() AuthorizedOpt {
	return func(a *Authorized) error {
		a.c.Set("payload", a.payload)
		return nil
	}
}

func SetUserData() AuthorizedOpt {
	return func(a *Authorized) error {
		user, err := a.mid.usersSvc.GetUser(a.payload.UserId)
		if err != nil {
			return err
		}

		a.c.Set("user", user)
		return nil
	}
}

func (m *Middleware) OnlyAuthorized(opts ...AuthorizedOpt) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload, err := m.authorized(c)
			if err != nil {
				utils.ClearCookies(c)
				c.Redirect(http.StatusFound, "/signin")
				return nil
			}

			a := &Authorized{
				c:       c,
				mid:     m,
				payload: payload,
			}

			for _, o := range opts {
				err := o(a)
				if err != nil {
					utils.ClearCookies(c)
					c.Redirect(http.StatusFound, "/signin")
					return nil
				}
			}

			return next(c)
		}
	}
}

func (m *Middleware) authorized(c echo.Context) (mytoken.Payload, error) {
	cookie, err := c.Cookie("accessToken")
	if err != nil {
		return mytoken.Payload{}, err
	}

	tokenString := cookie.Value
	if tokenString == "" {
		return mytoken.Payload{}, errors.New("invalid or empty token string")
	}

	claims, err := mytoken.ParseToken(mytoken.Access, m.cfg.Jwt.SecretKey, tokenString)
	if err != nil {
		return mytoken.Payload{}, err
	}

	has, err := m.usersSvc.HasAccessToken(claims.Payload.UserId, tokenString)
	if err != nil {
		return mytoken.Payload{}, err
	}
	if !has {
		return mytoken.Payload{}, errors.New("invalid or empty token string")
	}

	return claims.Payload, nil
}
