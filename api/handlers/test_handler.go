package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	tagsRepo  *repositories.TagsRepo
	codesRepo *repositories.CodesRepo
}

func NewTestHandler(
	tagsRepo *repositories.TagsRepo,
	codesRepo *repositories.CodesRepo,
) *TestHandler {
	return &TestHandler{
		tagsRepo:  tagsRepo,
		codesRepo: codesRepo,
	}
}

func (h *TestHandler) FindAllTags(c echo.Context) error {
	tags, err := h.tagsRepo.FindAllTags()
	if err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *TestHandler) FindAllCodes(c echo.Context) error {
	codes, err := h.codesRepo.FindAllCodes()
	if err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.JSON(http.StatusOK, codes)
}

func (h *TestHandler) FindOneCodeById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	code, err := h.codesRepo.FindOneCodeById(id)
	if err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.JSON(http.StatusOK, code)
}

func (h *TestHandler) CreateCode(c echo.Context) error {
	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadGateway, err.Error())
	}

	if err := utils.Validate(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	fmt.Println("req: ", req)

	if err := h.codesRepo.CreateCode(req); err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *TestHandler) UpdateCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := utils.Validate(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := h.codesRepo.UpdateCode(id, req); err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.NoContent(http.StatusOK)
}

func (h *TestHandler) DeleteCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := h.codesRepo.DeleteCode(id); err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.NoContent(http.StatusOK)
}
