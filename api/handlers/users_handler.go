package handlers

import (
	"errors"
	"fmt"
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
			ErrPayloadNotFound,
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

func (h *UsersHandler) UpdateUser(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(c, components.Error, ErrPayloadNotFound)
	}

	var req types.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	if err := utils.Validate(&req); err != nil {
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}
	files, ok := form.File["avatar"]
	if !ok {
		return utils.Render(c, components.Error("no \"avatar\" field"), http.StatusBadRequest)
	}

	if err = h.usersSvc.UpdateUser(payload.UserId, files[0], req); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return nil
}

func (h *UsersHandler) ToggleFollow(c echo.Context) error {
	payload, ok := c.Get("payload").(mytoken.Payload)
	if !ok {
		return utils.RenderError(c, components.Error, ErrPayloadNotFound)
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
			ErrPayloadNotFound,
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

	oldAccessCookie, _ := c.Cookie("accessToken")
	oldRefreshCookie, _ := c.Cookie("refreshToken")
	fmt.Printf("old cookies %s | %s\n", oldAccessCookie.Value, oldRefreshCookie.Value)

	utils.SetTokensCookies(c, token, h.cfg.Jwt)
	fmt.Printf("set new cookies to %+v\n", token)

	// test SetTokensCookies function
	accessCookie, _ := c.Cookie("accessToken")
	refreshCookie, _ := c.Cookie("refreshToken")
	if accessCookie.Value != token.AccessToken {
		return utils.RenderError(
			c,
			components.Error,
			errors.New(fmt.Sprintf(
				"accessCookie not equal: cookie=%s | expect=%s",
				accessCookie.Value,
				token.AccessToken,
			)),
		)
	}
	if refreshCookie.Value != token.RefreshToken {
		return utils.RenderError(
			c,
			components.Error,
			errors.New(fmt.Sprintf(
				"refreshCookie not equal: cookie=%s | expect=%s",
				refreshCookie.Value,
				token.RefreshToken,
			)),
		)
	}

	c.Response().Header().Set("HX-Trigger-After-Settle", "ready")

	return nil
}
