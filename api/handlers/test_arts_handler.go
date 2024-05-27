package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/pkg/httperror"
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
	arts, err := h.artsRepo.FindManyArts(1)
	if err != nil {
		_, status := httperror.Extract(err)
		return c.String(status, err.Error())
	}

	return c.JSON(http.StatusOK, arts)
}

func (h *testArtsHandler) FindOneArt(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_, status := httperror.Extract(err)
		return c.String(status, err.Error())
	}

	art, err := h.artsRepo.FindOneArt(id)
	if err != nil {
		_, status := httperror.Extract(err)
		return c.String(status, err.Error())
	}

	return c.JSON(http.StatusOK, art)
}
