package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
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
		return utils.Render(
			c,
			components.Error(err.Error()),
			http.StatusBadRequest,
		)
	}
	if err := utils.Validate(&req); err != nil {
		return utils.Render(
			c,
			components.Error(err.Error()),
			http.StatusBadRequest,
		)
	}

	arts, err := h.artsSvc.FindManyArts(req)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, components.Error(msg), status)
	}

	time.Sleep(300 * time.Millisecond)
	return utils.Render(c, components.ManyArts(arts), http.StatusOK)
}

func (h *ArtsHandler) ToggleStar(c echo.Context) error {
	errStatus := http.StatusInternalServerError
	errMsg := http.StatusText(errStatus)

	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.Render(c, components.Error(errMsg), errStatus)
	}

	artId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.Render(c, components.Error(errMsg), errStatus)
	}

	isStarred, err := h.artsSvc.ToggleStar(payload.UserId, artId)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, components.Error(msg), status)
	}

	return utils.Render(c, components.StarButton(artId, isStarred), http.StatusOK)
}
