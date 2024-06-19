package handlers

import (
	"net/http"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
)

type TestTagsHandler struct {
	tagsRepo *repositories.TagsRepo
}

func NewTestTagsHandler(tagsRepo *repositories.TagsRepo) *TestTagsHandler {
	return &TestTagsHandler{
		tagsRepo: tagsRepo,
	}
}

func (h *TestTagsHandler) FindAllTags(c echo.Context) error {
	tags, err := h.tagsRepo.FindAllTags()
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, tags)
}
