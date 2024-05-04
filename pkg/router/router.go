package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(e *echo.Echo) *Router {
	return &Router{
		e: e,
	}
}

func (r *Router) TestRouter() {
	r.e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})
}
