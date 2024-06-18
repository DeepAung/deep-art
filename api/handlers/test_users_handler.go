package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/repositories"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/labstack/echo/v4"
)

type TestUsersHandler struct {
	usersRepo *repositories.UsersRepo
}

func NewTestUsersHandler(usersRepo *repositories.UsersRepo) *TestUsersHandler {
	return &TestUsersHandler{
		usersRepo: usersRepo,
	}
}

func (h *TestUsersHandler) GetCreator(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	creator, err := h.usersRepo.FindOneCreatorById(id)
	if err != nil {
		msg, status := httperror.Extract(err)
		return c.JSON(status, msg)
	}

	return c.JSON(http.StatusOK, creator)
}
