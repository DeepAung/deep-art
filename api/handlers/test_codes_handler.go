package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
)

type testCodesHandler struct {
	codesRepo *repositories.CodesRepo
}

func NewTestCodesHandler(codesRepo *repositories.CodesRepo) *testCodesHandler {
	return &testCodesHandler{
		codesRepo: codesRepo,
	}
}

func (h *testCodesHandler) FindAllCodes(c echo.Context) error {
	codes, err := h.codesRepo.FindAllCodes()
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, codes)
}

func (h *testCodesHandler) FindOneCodeById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	code, err := h.codesRepo.FindOneCodeById(id)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, code)
}

func (h *testCodesHandler) CreateCode(c echo.Context) error {
	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}

	if err := utils.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.codesRepo.CreateCode(req); err != nil {
		return utils.JSONError(c, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *testCodesHandler) UpdateCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var req types.CodeReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := utils.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.codesRepo.UpdateCode(id, req); err != nil {
		return utils.JSONError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *testCodesHandler) DeleteCode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.codesRepo.DeleteCode(id); err != nil {
		return utils.JSONError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *testCodesHandler) UseCode(c echo.Context) error {
	userId, err := strconv.Atoi(c.FormValue("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	codeId, err := strconv.Atoi(c.FormValue("codeId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return h.codesRepo.UseCode(userId, codeId)
}
