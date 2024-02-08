package tagsHandler

import (
	"strconv"

	"github.com/DeepAung/deep-art/modules/tags"
	"github.com/DeepAung/deep-art/modules/tags/tagsUsecase"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/gofiber/fiber/v2"
)

const (
	_ = response.TraceId(
		"users-" + string('0'+iota/100%10) + string('0'+iota/10%10) + string('0'+iota/1%10),
	)

	getTagsErr
	createTagErr
	updateTagErr
	deleteTagErr
)

type ITagsHandler interface {
	GetTags(c *fiber.Ctx) error
	CreateTag(c *fiber.Ctx) error
	UpdateTag(c *fiber.Ctx) error
	DeleteTag(c *fiber.Ctx) error
}

type tagsHandler struct {
	tagsUsecase tagsUsecase.ITagsUsecase
}

func NewTagsHandler(tagsUsecase tagsUsecase.ITagsUsecase) ITagsHandler {
	return &tagsHandler{
		tagsUsecase: tagsUsecase,
	}
}

func (h *tagsHandler) GetTags(c *fiber.Ctx) error {
	tags, err := h.tagsUsecase.GetTags()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, getTagsErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, tags)
}

func (h *tagsHandler) CreateTag(c *fiber.Ctx) error {
	req := new(tags.TagReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, createTagErr, err.Error())
	}

	err := h.tagsUsecase.CreateTag(req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, createTagErr, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, nil)
}

func (h *tagsHandler) UpdateTag(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, updateTagErr, "invalid tag id")
	}

	req := new(tags.TagReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, updateTagErr, err.Error())
	}

	err = h.tagsUsecase.UpdateTag(req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, updateTagErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil)
}

func (h *tagsHandler) DeleteTag(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, updateTagErr, "invalid tag id")
	}

	err = h.tagsUsecase.DeleteTag(id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, deleteTagErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil)
}
