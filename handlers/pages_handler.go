package handlers

import (
	"context"

	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
)

type PagesHandler struct{}

func NewPagesHandler() *PagesHandler {
	return &PagesHandler{}
}

func (h *PagesHandler) Welcome(c echo.Context) error {
	return pages.Welcome().Render(context.Background(), c.Response())
}
