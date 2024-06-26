package handlers

import (
	"net/http"

	"github.com/DeepAung/deep-art/pkg/storer"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/labstack/echo/v4"
)

type TestFilesHandler struct {
	storer storer.Storer
}

func NewTestFilesHandler(storer storer.Storer) *TestFilesHandler {
	return &TestFilesHandler{
		storer: storer,
	}
}

func (h *TestFilesHandler) UploadFiles(c echo.Context) error {
	dir := c.FormValue("dir")
	form, err := c.MultipartForm()
	if err != nil {
		return utils.JSONError(c, err)
	}

	files, ok := form.File["files"]
	if !ok {
		return c.JSON(http.StatusBadRequest, "no files field")
	}

	res, err := h.storer.UploadFiles(files, dir)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *TestFilesHandler) DeleteFiles(c echo.Context) error {
	var req struct {
		Dests []string `json:"dests"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.storer.DeleteFiles(req.Dests); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
