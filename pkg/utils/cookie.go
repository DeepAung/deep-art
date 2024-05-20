package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetCookie(c echo.Context, name, value string, maxAge time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(maxAge),
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
