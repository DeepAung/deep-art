package handlers

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
)

type CodesHandler struct {
	codesSvc *services.CodesSvc
	cfg      *config.Config
}

func NewCodesHandler(codesSvc *services.CodesSvc, cfg *config.Config) *CodesHandler {
	return &CodesHandler{
		codesSvc: codesSvc,
		cfg:      cfg,
	}
}

func (h *CodesHandler) UseCode(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		// TODO: toast error
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		return utils.Render(c, components.Error(msg), status)
	}

	name := c.FormValue("name")
	if name == "" {
		// TODO: toast error
		return utils.Render(c, components.Error("code should not be empty"), http.StatusBadRequest)
	}

	err := h.codesSvc.UseCode(payload.UserId, name)
	if err != nil {
		msg, status := httperror.Extract(err)
		return utils.Render(c, components.Error(msg), status)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return nil
}
