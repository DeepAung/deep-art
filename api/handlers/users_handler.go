package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/mytoken"
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

	if err := utils.Validate(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	passport, err := h.usersSvc.SignIn(req.Email, req.Password)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	utils.SetTokensCookies(c, passport.Token, h.cfg.Jwt)

	c.Response().Header().Add("HX-Redirect", "/home")
	return nil
}

func (h *UsersHandler) SignUp(c echo.Context) error {
	var req types.SignUpReq
	if err := c.Bind(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	if err := utils.Validate(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	_, err := h.usersSvc.SignUp(req)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Redirect", "/signin")
	return nil
}

func (h *UsersHandler) SignOut(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			errors.New("payload from middleware not found"),
		)
	}

	cookie, err := c.Cookie("tokenId")
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}
	tokenId, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	err = h.usersSvc.SignOut(payload.UserId, tokenId)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	utils.ClearTokensCookies(c)
	c.Response().Header().Add("HX-Redirect", "/signin")
	return nil
}

func (h *UsersHandler) ToggleFollow(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			errors.New("payload from middleware not found"),
		)
	}

	creatorId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	isFollowing, err := h.usersSvc.ToggleFollow(payload.UserId, creatorId)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	return utils.Render(c, components.FollowButton(creatorId, isFollowing), http.StatusOK)
}

func (h *UsersHandler) UpdateTokens(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(
			c,
			components.Error,
			errors.New("payload from middleware not found"),
		)
	}

	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	token, err := h.usersSvc.UpdateTokens(payload.UserId, refreshTokenCookie.Value)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	utils.SetTokensCookies(c, token, h.cfg.Jwt)

	return nil
}
