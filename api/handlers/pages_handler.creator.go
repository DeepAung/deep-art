package handlers

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
)

func (h *PagesHandler) CreatorHomePage(c echo.Context) error {
	user, ok := c.Get("user").(types.User)
	if !ok {
		return utils.RenderError(c, pages.Error, ErrUserDataNotFound)
	}

	return utils.Render(c, pages.CreatorHome(user), http.StatusOK)
}

func (h *PagesHandler) CreatorCreateArt(c echo.Context) error {
	user, ok := c.Get("user").(types.User)
	if !ok {
		return utils.RenderError(c, pages.Error, ErrUserDataNotFound)
	}

	return utils.Render(c, pages.CreatorCreateArt(user), http.StatusOK)
}
