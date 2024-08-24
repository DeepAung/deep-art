package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
)

type TagsHandler struct {
	tagsSvc *services.TagsSvc
}

func NewTagsHandler(tagsSvc *services.TagsSvc) *TagsHandler {
	return &TagsHandler{
		tagsSvc: tagsSvc,
	}
}

func (h *TagsHandler) TagsFilter(c echo.Context) error {
	tags, err := h.tagsSvc.GetTags()
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.HomeTagsFilter(tags), http.StatusOK)
}

func (h *TagsHandler) TagsOptions(c echo.Context) error {
	tags, err := h.tagsSvc.GetTags()
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.TagsOptions(tags), http.StatusOK)
}

func (h *TagsHandler) GetTags(c echo.Context) error {
	tags, err := h.tagsSvc.GetTags()
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}
	return utils.Render(c, components.Tags(tags), http.StatusOK)
}

func (h *TagsHandler) CreateTag(c echo.Context) error {
	var req types.TagReq
	if err := c.Bind(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	tag, err := h.tagsSvc.CreateTag(req.Name)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#create-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.Tag(tag), http.StatusCreated)
}

func (h *TagsHandler) UpdateTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	var req types.TagReq
	if err := c.Bind(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	tag, err := h.tagsSvc.UpdateTag(id, req.Name)
	if err != nil {
		c.Response().Header().Add("HX-Retarget", "#update-tag-error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.Tag(tag), http.StatusOK)
}

func (h *TagsHandler) DeleteTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	if err := h.tagsSvc.DeleteTag(id); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return nil
}
