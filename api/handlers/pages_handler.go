package handlers

import (
	"context"
	"net/http"

	"github.com/DeepAung/deep-art/pkg/mytoken"
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

func (h *PagesHandler) Home(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return pages.Error(http.StatusText(http.StatusInternalServerError)).Render(context.Background(), c.Response())
	}

	return pages.Home(payload).Render(context.Background(), c.Response())
}

func (h *PagesHandler) SignIn(c echo.Context) error {
	return pages.SignIn().Render(context.Background(), c.Response())
}

func (h *PagesHandler) SignUp(c echo.Context) error {
	return pages.SignUp().Render(context.Background(), c.Response())
}
