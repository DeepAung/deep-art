package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/DeepAung/deep-art/pkg/httperror"
	"github.com/DeepAung/deep-art/pkg/utils"
	"github.com/DeepAung/deep-art/views/components"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (h *UsersHandler) OAuthHandler(c echo.Context) error {
	addProviderParamInQuery(c)

	if gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		return h.oauthCallbackSignin(c, gothUser, "")
	}

	utils.SetCookie(c, "redirect_to", c.QueryParam("redirect_to"), 0)
	utils.SetCookie(c, "callback_func", c.QueryParam("callback_func"), 0)

	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

func (h *UsersHandler) OAuthCallback(c echo.Context) error {
	addProviderParamInQuery(c)

	gothUser, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return utils.RenderError(c, components.Error, err)
	}

	callbackFunc, err := utils.GetCookie(c, "callback_func", "")
	if err != nil {
		return err
	}
	redirectTo, err := utils.GetCookie(c, "redirect_to", "")
	if err != nil {
		return err
	}
	utils.DeleteCookie(c, "redirect_to")
	utils.DeleteCookie(c, "callback_func")

	switch callbackFunc {
	case "signin":
		return h.oauthCallbackSignin(c, gothUser, redirectTo)
	case "signup":
		return h.oauthCallbackSignup(c, gothUser, redirectTo)
	case "connect":
		return h.oauthCallbackConnect(c, gothUser)
	case "disconnect":
		return h.oauthCallbackDisconnect(c, gothUser)
	default:
		return utils.RenderError(c, components.Error, errors.New("invalid callback_func"))
	}
}

// TODO:
// find user with gothUser.Email check if user is connect OAuth to this provider
// redirect back
func (h *UsersHandler) oauthCallbackSignin(
	c echo.Context,
	gothUser goth.User,
	redirectTo string,
) error {
	if redirectTo == "" {
		redirectTo = "/home"
	}

	passport, err := h.usersSvc.OAuthSignin(gothUser, redirectTo)
	if err != nil {
		msg, status := httperror.Extract(err)
		if status >= 500 {
			slog.Error(err.Error())
		}
		return utils.Render(c, components.ErrorWithBackBtn(msg, "/signin"), status)
	}

	utils.SetTokensCookies(
		c,
		passport.Token.Id,
		passport.Token.AccessToken,
		passport.Token.RefreshToken,
		h.cfg.Jwt,
	)
	utils.SetCookie(c, "provider", gothUser.Provider, 0)

	return c.Redirect(http.StatusFound, redirectTo)
}

// TODO:
// create new user with this gothUser info
// if there is already the same email (show error / connect&signin)
func (h *UsersHandler) oauthCallbackSignup(
	c echo.Context,
	gothUser goth.User,
	redirectTo string,
) error {
	if redirectTo == "" {
		redirectTo = "/signin"
	}

	if _, err := h.usersSvc.OAuthSignup(gothUser, redirectTo); err != nil {
		msg, status := httperror.Extract(err)
		if status >= 500 {
			slog.Error(err.Error())
		}
		return utils.Render(c, components.ErrorWithBackBtn(msg, "/signup"), status)
	}

	return c.Redirect(http.StatusFound, redirectTo)
}

// TODO:
func (h *UsersHandler) oauthCallbackConnect(
	c echo.Context,
	gothUser goth.User,
) error {
	return nil
}

// TODO:
func (h *UsersHandler) oauthCallbackDisconnect(
	c echo.Context,
	gothUser goth.User,
) error {
	return nil
}

func addProviderParamInQuery(c echo.Context) {
	req := c.Request()
	q := req.URL.Query()
	q.Add("provider", c.Param("provider"))
	req.URL.RawQuery = q.Encode()
	c.SetRequest(req)
}
