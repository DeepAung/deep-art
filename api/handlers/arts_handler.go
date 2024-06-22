package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
)

type ArtsHandler struct {
	artsSvc *services.ArtsSvc
	cfg     *config.Config
}

func NewArtsHandler(artsSvc *services.ArtsSvc, cfg *config.Config) *ArtsHandler {
	return &ArtsHandler{
		artsSvc: artsSvc,
		cfg:     cfg,
	}
}

func (h *ArtsHandler) FindManyArts(c echo.Context) error {
	var req types.ManyArtsReq
	if err := c.Bind(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	arts, err := h.artsSvc.FindManyArts(req)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	time.Sleep(300 * time.Millisecond)
	return utils.Render(c, components.ManyArts(arts), http.StatusOK)
}

func (h *ArtsHandler) BuyArt(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			errors.New("payload from middleware not found"),
		)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	if err := h.artsSvc.BuyArt(payload.UserId, artId, price); err != nil {
		return utils.RenderError(c, components.ErrToast, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return nil
}

func (h *ArtsHandler) ToggleStar(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			ErrPayloadNotFound,
		)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	isStarred, err := h.artsSvc.ToggleStar(payload.UserId, artId)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.StarButton(artId, isStarred), http.StatusOK)
}
