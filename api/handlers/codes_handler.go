package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
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
		return utils.RenderError(
			c,
			components.Error,
			ErrPayloadNotFound,
		)
	}

	name := c.FormValue("name")
	if name == "" {
		return utils.Render(c, components.Error("code should not be empty"), http.StatusBadRequest)
	}

	err := h.codesSvc.UseCode(payload.UserId, name)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
	return nil
}

func (h *CodesHandler) CreateCode(c echo.Context) error {
	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	code, err := h.codesSvc.CreateCode(req)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.Code(code), http.StatusCreated)
}

func (h *CodesHandler) GetCodes(c echo.Context) error {
	codes, err := h.codesSvc.GetCodes()
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.Codes(codes), http.StatusOK)
}

func (h *CodesHandler) UpdateCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	code, err := h.codesSvc.UpdateCode(id, req)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-code-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.Code(code), http.StatusOK)
}

func (h *CodesHandler) DeleteCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	if err := h.codesSvc.DeleteCode(id); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return nil
}
