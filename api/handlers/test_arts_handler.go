package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
)

type testArtsHandler struct {
	artsRepo *repositories.ArtsRepo
}

func NewTestArtsHandler(artsRepo *repositories.ArtsRepo) *testArtsHandler {
	return &testArtsHandler{
		artsRepo: artsRepo,
	}
}

func (h *testArtsHandler) FindManyArts(c echo.Context) error {
	var req types.ManyArtsReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	arts, err := h.artsRepo.FindManyArts(req)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, arts)
}

func (h *testArtsHandler) FindOneArt(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	art, err := h.artsRepo.FindOneArt(id)
	if err != nil {
		return utils.JSONError(c, err)
	}

	return c.JSON(http.StatusOK, art)
}
