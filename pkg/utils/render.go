package utils

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component, status int) error {
	c.Response().Status = status
	err := component.Render(c.Request().Context(), c.Response())
	if err == nil {
		return nil
	}

	status = http.StatusInternalServerError
	return c.String(status, http.StatusText(status))
}
