package utils

import (
	"log/slog"
	"net/http"

	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func JSONError(c echo.Context, err error) error {
	msg, status := httperror.Extract(err)
	if status >= 500 {
		slog.Error(err.Error())
	}
	return c.JSON(status, msg)
}

func RenderError(c echo.Context, errorComponent func(msg string) templ.Component, err error) error {
	msg, status := httperror.Extract(err)
	if status >= 500 {
		slog.Error(err.Error())
	}
	return Render(c, errorComponent(msg), status)
}

func Render(c echo.Context, component templ.Component, status int) error {
	c.Response().Status = status
	err := component.Render(c.Request().Context(), c.Response().Writer)
	if err == nil {
		return nil
	}

	status = http.StatusInternalServerError
	return c.String(status, http.StatusText(status))
}
