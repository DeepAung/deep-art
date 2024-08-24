package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

func (h *PagesHandler) CreatorArtDetail(c echo.Context) error {
	user, ok := c.Get("user").(types.User)
	if !ok {
		return utils.RenderError(c, pages.Error, ErrUserDataNotFound)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Render(c, pages.Error("Page Not Found"), http.StatusNotFound)
	}

	art, err := h.artsSvc.FindOneArt(artId)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return utils.RenderError(c, pages.Error, err)
	}

	tags, err := h.tagsSvc.GetTags()
	if err != nil {
		fmt.Println("err2: ", err.Error())
		return utils.RenderError(c, pages.Error, err)
	}

	return utils.Render(c, pages.CreatorArtDetail(user, art, tags), http.StatusOK)
}

func (h *PagesHandler) CreatorProfile(c echo.Context) error {
	me, ok := c.Get("user").(types.User)
	if !ok {
		return utils.RenderError(c, pages.Error, ErrUserDataNotFound)
	}

	creatorId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Render(c, pages.Error("Page Not Found"), http.StatusNotFound)
	}

	creator, err := h.usersSvc.GetCreator(creatorId)
	if err != nil {
		return utils.RenderError(c, pages.Error, err)
	}

	return utils.Render(c, pages.CreatorProfile(me, creator), http.StatusOK)
}
