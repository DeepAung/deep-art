package handlers

import (
	"log/slog"
	"net/http"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/pkg/httperror"
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
	tags, err := h.tagsSvc.FindAllTags()
	if err != nil {
		slog.Error(err.Error())
		msg, status := httperror.Extract(err)
		return utils.Render(c, components.Error(msg), status)
	}

	return utils.Render(c, components.HomeTagsFilter(tags), http.StatusOK)
}
