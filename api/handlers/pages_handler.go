package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
)

type PagesHandler struct {
	usersSvc *services.UsersSvc
	artsSvc  *services.ArtsSvc
}

func NewPagesHandler(usersSvc *services.UsersSvc, artsSvc *services.ArtsSvc) *PagesHandler {
	return &PagesHandler{
		usersSvc: usersSvc,
		artsSvc:  artsSvc,
	}
}

func (h *PagesHandler) Welcome(c echo.Context) error {
	return utils.Render(c, pages.Welcome(), http.StatusOK)
}

func (h *PagesHandler) SignIn(c echo.Context) error {
	return utils.Render(c, pages.SignIn(), http.StatusOK)
}

func (h *PagesHandler) SignUp(c echo.Context) error {
	return utils.Render(c, pages.SignUp(), http.StatusOK)
}

func (h *PagesHandler) Home(c echo.Context) error {
	user, ok := c.Get("user").(types.User)
	if !ok {
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		return utils.Render(c, pages.Error(msg), status)
	}

	return utils.Render(c, pages.Home(user), http.StatusOK)
}

func (h *PagesHandler) ArtDetail(c echo.Context) error {
	user, ok := c.Get("user").(types.User)
	if !ok {
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		return utils.Render(c, pages.Error(msg), status)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Render(c, pages.Error("Page Not Found"), http.StatusNotFound)
	}

	art, err := h.artsSvc.FindOneArt(artId)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, pages.Error(msg), status)
	}

	isFollowing, err := h.usersSvc.IsFollowing(user.Id, art.Creator.Id)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, pages.Error(msg), status)
	}

	isStarred, err := h.artsSvc.IsStarred(user.Id, artId)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, pages.Error(msg), status)
	}

	return utils.Render(c, pages.ArtDetail(user, art, isFollowing, isStarred), http.StatusOK)
}
