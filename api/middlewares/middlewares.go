package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/prettylog"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	emptyTokenErr   = httperror.New("invalid or empty token string", http.StatusBadRequest)
	invalidTokenErr = httperror.New("invalid or empty token string", http.StatusBadRequest)
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

func clearCookieAndRedirect(c echo.Context) error {
	utils.ClearTokensCookies(c)

	myUrl, err := url.Parse("/signin")
	if err != nil {
		return err
	}

	q := myUrl.Query()
	q.Add("redirect_to", c.Request().URL.String())
	myUrl.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, myUrl.String())

	return nil
}

func (m *Middleware) OnlyAuthorized(opts ...AuthorizedOpt) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, payload, err := m.jwtAccessToken(c)
			if err != nil {
				fmt.Println("firsttry err: ", err.Error())
				if !errors.Is(err, http.ErrNoCookie) && !errors.Is(err, emptyTokenErr) {
					return clearCookieAndRedirect(c)
				}

				payload, err = m.tryUpdateToken(c)
				if err != nil {
					fmt.Println("tryupdatetoken err: ", err.Error())
					return clearCookieAndRedirect(c)
				}
			}

			a := &Authorized{
				c:       c,
				mid:     m,
				payload: payload,
			}

			for _, o := range opts {
				err := o(a)
				if err != nil {
					return clearCookieAndRedirect(c)
				}
			}

			return next(c)
		}
	}
}

func (m *Middleware) JwtRefreshToken(opts ...AuthorizedOpt) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, payload, err := m.jwtRefreshToken(c)
			if err != nil {
				return clearCookieAndRedirect(c)
			}

			a := &Authorized{
				c:       c,
				mid:     m,
				payload: payload,
			}

			for _, o := range opts {
				err := o(a)
				if err != nil {
					return clearCookieAndRedirect(c)
				}
			}

			return next(c)
		}
	}
}

func (m *Middleware) jwtAccessToken(c echo.Context) (string, mytoken.Payload, error) {
	cookie, err := c.Cookie("accessToken")
	if err != nil {
		fmt.Println("err no cookie")
		return "", mytoken.Payload{}, http.ErrNoCookie
	}

	tokenString := cookie.Value
	if tokenString == "" {
		fmt.Println("err no token")
		return "", mytoken.Payload{}, emptyTokenErr
	}

	claims, err := mytoken.ParseToken(mytoken.Access, m.cfg.Jwt.SecretKey, tokenString)
	if err != nil {
		fmt.Println("err parse token")
		return "", mytoken.Payload{}, err
	}

	has, err := m.usersSvc.HasAccessToken(claims.Payload.UserId, tokenString)
	if err != nil {
		fmt.Println("has access err")
		return "", mytoken.Payload{}, err
	}
	if !has {
		fmt.Println("err !has")
		return "", mytoken.Payload{}, invalidTokenErr
	}

	return tokenString, claims.Payload, nil
}

func (m *Middleware) jwtRefreshToken(c echo.Context) (string, mytoken.Payload, error) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return "", mytoken.Payload{}, http.ErrNoCookie
	}

	tokenString := cookie.Value
	if tokenString == "" {
		return "", mytoken.Payload{}, emptyTokenErr
	}

	claims, err := mytoken.ParseToken(mytoken.Refresh, m.cfg.Jwt.SecretKey, tokenString)
	if err != nil {
		return "", mytoken.Payload{}, err
	}

	has, err := m.usersSvc.HasRefreshToken(claims.Payload.UserId, tokenString)
	if err != nil {
		return "", mytoken.Payload{}, err
	}
	if !has {
		return "", mytoken.Payload{}, invalidTokenErr
	}

	return tokenString, claims.Payload, nil
}

func (m *Middleware) tryUpdateToken(c echo.Context) (mytoken.Payload, error) {
	refreshToken, payload, err := m.jwtRefreshToken(c)
	if err != nil {
		return mytoken.Payload{}, err
	}

	token, err := m.usersSvc.UpdateTokens(payload.UserId, refreshToken)
	if err != nil {
		return mytoken.Payload{}, err
	}

	utils.SetTokensCookies(c, token, m.cfg.Jwt)
	return payload, nil
}
