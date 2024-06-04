package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/pages"
	"github.com/labstack/echo/v4"
)

type PagesHandler struct {
	ArtsSvc *services.ArtsSvc
}

func NewPagesHandler(ArtsSvc *services.ArtsSvc) *PagesHandler {
	return &PagesHandler{
		ArtsSvc: ArtsSvc,
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
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		return utils.Render(c, pages.Error(msg), status)
	}

	return utils.Render(c, pages.Home(payload), http.StatusOK)
}

func (h *PagesHandler) ArtDetail(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		return utils.Render(c, pages.Error(msg), status)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Render(c, pages.Error("Page Not Found"), http.StatusNotFound)
	}

	art, err := h.ArtsSvc.FindOneArt(id)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, pages.Error(msg), status)
	}

	return utils.Render(c, pages.ArtDetail(payload, art), http.StatusOK)
}
