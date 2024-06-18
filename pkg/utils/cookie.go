package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DeepAung/deep-art/api/types"
	"github.com/DeepAung/deep-art/pkg/config"
	"github.com/labstack/echo/v4"
)

func SetCookie(c echo.Context, name, value string, maxAge time.Duration) {
	var expires time.Time
	if maxAge != 0 {
		expires = time.Now().Add(maxAge)
	}

	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(maxAge.Seconds()),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func DeleteCookie(c echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-7 * 24 * time.Hour),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func SetTokensCookies(c echo.Context, token types.Token, cfg *config.JwtConfig) {
	SetCookie(c, "accessToken", token.AccessToken, cfg.AccessExpires)
	SetCookie(c, "refreshToken", token.RefreshToken, cfg.RefreshExpires)
	SetCookie(c, "tokenId", strconv.Itoa(token.Id), 0)
}

func ClearTokensCookies(c echo.Context) {
	DeleteCookie(c, "accessToken")
	DeleteCookie(c, "refreshToken")
	DeleteCookie(c, "tokenId")
}
