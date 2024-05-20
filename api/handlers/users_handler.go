package handlers

import (
	"log/slog"
	"net/http"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
)

type UsersHandler struct {
	usersSvc *services.UsersSvc
	cfg      *config.Config
}

func NewUsersHandler(usersSvc *services.UsersSvc, cfg *config.Config) *UsersHandler {
	return &UsersHandler{
		usersSvc: usersSvc,
		cfg:      cfg,
	}
}

func (h *UsersHandler) SignIn(c echo.Context) error {
	var req types.SignInReq
	if err := c.Bind(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	if err := utils.Validate(req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	passport, err := h.usersSvc.SignIn(req.Email, req.Password)
	if err != nil {
		msg, status := httperror.Extract(err)
		slog.Error(err.Error())
		return utils.Render(c, components.Error(msg), status)
	}

	utils.SetCookie(c, "accessToken", passport.Token.AccessToken, h.cfg.Jwt.AccessExpires)
	utils.SetCookie(c, "refreshToken", passport.Token.RefreshToken, h.cfg.Jwt.RefreshExpires)

	c.Response().Header().Add("HX-Redirect", "/home")
	return nil
}

func (h *UsersHandler) SignUp(c echo.Context) error {
	var req types.SignUpReq
	if err := c.Bind(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	if err := utils.Validate(req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	_, err := h.usersSvc.SignUp(req)
	if err != nil {
		msg, status := httperror.Extract(err)
		slog.Error(err.Error())
		return utils.Render(c, components.Error(msg), status)
	}

	c.Response().Header().Add("HX-Redirect", "/signin")
	return nil
}

func (h *UsersHandler) SignOut(c echo.Context) error {
	h.usersSvc.SignOut(userId int, tokenId int)
}
