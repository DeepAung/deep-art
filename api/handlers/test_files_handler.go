package handlers

import (
	"io"
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

	files2 := make([]io.Reader, len(files))
	dests := make([]string, len(files))
	for i := range len(files) {
		f, err := files[i].Open()
		if err != nil {
			return err
		}
		defer f.Close()

		files2[i] = f
		dests[i] = utils.Join(dir, files[i].Filename)
	}

	res, err := h.storer.UploadFiles(files2, dests)
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
