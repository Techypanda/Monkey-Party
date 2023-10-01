package whoami

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var sessGet = session.Get

func WhoAmI(c echo.Context) error {
	sess, _ := sessGet("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	if userId, ok := sess.Values["user"].(string); ok {
		return c.JSON(http.StatusOK, map[string]string{
			"userId": userId,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"message": "no session",
	})
}
