package router

import (
	"github.com/labstack/echo/v4"
	"techytechster.com/monkeyparty/internal/handlers/ping"
	"techytechster.com/monkeyparty/internal/handlers/session"
	"techytechster.com/monkeyparty/internal/handlers/whoami"
)

var (
	PingHandler        = ping.Ping
	SessionHandler     = session.AcquireSession
	AcquireUserHandler = whoami.WhoAmI
)

func Router(e *echo.Echo) {
	e.GET("/ping", PingHandler)
	e.GET("/session", SessionHandler)
	e.GET("/whoami", AcquireUserHandler)
}
