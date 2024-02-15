package filesHandler

import (
	"github.com/DeepAung/deep-art/modules/files"
	"github.com/DeepAung/deep-art/pkg/mystorage"
	"github.com/DeepAung/deep-art/pkg/response"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadFilesErr response.TraceId = "files-001"
	deleteFilesErr response.TraceId = "files-002"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFiles(c *fiber.Ctx) error
}

type filesHandler struct {
	s mystorage.IStorage
}

func NewFilesHandler(s mystorage.IStorage) IFilesHandler {
	return &filesHandler{
		s: s,
	}
}
func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, uploadFilesErr, err.Error())
	}

	dir := c.FormValue("dir")

	files, ok := form.File["files"]
	if !ok {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			uploadFilesErr,
			"\"files\" field not found",
		)
	}

	results, err := h.s.UploadFiles(files, dir)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, uploadFilesErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, results)
}

func (h *filesHandler) DeleteFiles(c *fiber.Ctx) error {
	req := new(files.DeleteFilesReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, deleteFilesErr, err.Error())
	}

	err := h.s.DeleteFiles(req.Destinations)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, deleteFilesErr, err.Error())
	}

	return response.Success(c, fiber.StatusOK, nil)
}
