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
		return h.oauthCallbackSignin(c, gothUser)
	}

	c.SetCookie(&http.Cookie{
		Name:  "redirect_to",
		Value: c.QueryParam("redirect_to"),
		Path:  "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:  "callback_func",
		Value: c.QueryParam("callback_func"),
		Path:  "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})

	// TODO: test redirect_to cookie
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

	switch callbackFunc {
	case "signin":
		return h.oauthCallbackSignin(c, gothUser)
	case "signup":
		return h.oauthCallbackSignup(c, gothUser)
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
func (h *UsersHandler) oauthCallbackSignin(c echo.Context, gothUser goth.User) error {
	redirectTo, err := utils.GetCookie(c, "redirect_to", "/home")
	if err != nil {
		return err
	}
	_, _ = redirectTo, gothUser
	return nil
}

// TODO:
// create new user with this gothUser info
// if there is already the same email (show error / connect&signin)
func (h *UsersHandler) oauthCallbackSignup(c echo.Context, gothUser goth.User) error {
	redirectTo, err := utils.GetCookie(c, "redirect_to", "/signin")
	if err != nil {
		return err
	}

	if _, err := h.usersSvc.OAuthSignup(gothUser, redirectTo); err != nil {
		msg, status := httperror.Extract(err)
		if status >= 500 {
			slog.Error(err.Error())
		}
		return utils.Render(c, components.ErrorWithBackBtn(msg, "/signup"), status)
	}

	return c.Redirect(http.StatusPermanentRedirect, redirectTo)
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

// TODO: maybe merge this function with Signout()
func (h *UsersHandler) OAuthSignout(c echo.Context) error {
	addProviderParamInQuery(c)
	gothic.Logout(c.Response(), c.Request())
	return nil
}

func addProviderParamInQuery(c echo.Context) {
	req := c.Request()
	q := req.URL.Query()
	q.Add("provider", c.Param("provider"))
	req.URL.RawQuery = q.Encode()
	c.SetRequest(req)
}
