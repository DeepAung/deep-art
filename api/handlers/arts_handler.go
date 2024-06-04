package handlers

import (
	"net/http"
	"time"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
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
		return err
		// return utils.Render(
		// 	c,
		// 	components.Error("fetching arts failed"),
		// 	http.StatusInternalServerError,
		// )
	}
	if err := utils.Validate(&req); err != nil {
		return err
	}

	arts, err := h.artsSvc.FindManyArts(req.Page)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, components.Error(msg), status)
	}

	time.Sleep(300 * time.Millisecond)
	return utils.Render(c, components.ManyArts(arts), http.StatusOK)
}
