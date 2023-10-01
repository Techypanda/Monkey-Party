package session

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var sessGet = session.Get

func AcquireSession(c echo.Context) error {
	sess, _ := sessGet("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	sess.Values["user"] = uuid.New().String()
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}
