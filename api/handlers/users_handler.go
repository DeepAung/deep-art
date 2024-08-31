package handlers

import (
	"net/http"
	"strconv"

	"github.com/DeepAung/deep-art/api/services"
	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/DeepAung/deep-art/pkg/mytoken"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
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
	if req.RedirectTo == "" {
		req.RedirectTo = "/home"
	}

	passport, err := h.usersSvc.SignIn(req.Email, req.Password)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	utils.SetTokensCookies(
		c,
		passport.Token.Id,
		passport.Token.AccessToken,
		passport.Token.RefreshToken,
		h.cfg.Jwt,
	)

	c.Response().Header().Add("HX-Redirect", req.RedirectTo)
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
	if req.RedirectTo == "" {
		req.RedirectTo = "/signin"
	}

	_, err := h.usersSvc.SignUp(req)
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Redirect", req.RedirectTo)
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

	// try to logout OAuth
	if provider, err := utils.GetCookie(c, "provider", ""); err == nil && provider != "" {
		// add provider cookie in query
		req := c.Request()
		q := req.URL.Query()
		q.Add("provider", provider)
		req.URL.RawQuery = q.Encode()
		c.SetRequest(req)

		utils.DeleteCookie(c, "provider")
		gothic.Logout(c.Response(), c.Request())
	}

	if err = h.usersSvc.SignOut(payload.UserId, tokenId); err != nil {
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
		return utils.Render(c, components.Error(err.Error()), http.StatusBadRequest)
	}
	files, ok := form.File["avatar"]
	if !ok {
		return utils.Render(c, components.Error("no \"avatar\" field"), http.StatusBadRequest)
	}

	if err = h.usersSvc.UpdateUser(payload.UserId, files[0], req); err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	c.Response().Header().Add("HX-Refresh", "true")
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

	utils.SetTokensCookies(c, token.Id, token.AccessToken, token.RefreshToken, h.cfg.Jwt)

	c.Response().Header().Set("HX-Trigger-After-Settle", "ready")

	return nil
}
