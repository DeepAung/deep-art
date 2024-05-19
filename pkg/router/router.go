package router

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/handlers"
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

func (r *Router) PagesRouter() {
	handler := handlers.NewPagesHandler()

	r.e.GET("/", handler.Welcome)
}

func (r *Router) TestRouter() {
	r.e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test route")
	})
}
