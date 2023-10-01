package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"techytechster.com/monkeyparty/internal/router"
	"techytechster.com/monkeyparty/internal/utils"
)

var e = echo.New()
var fatal = e.Logger.Fatal
var start = e.Start

func logValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	logger := zerolog.New(os.Stdout)
	logger.Info().
		Str("URI", v.URI).
		Int("status", v.Status).
		Msg("request")
	return nil
}

const SECRET_LENGTH = 4128
const BURST_TPS = 10
const RATE_TPS = 5

func identifierExtractor(ctx echo.Context) (string, error) {
	id := ctx.RealIP()
	return id, nil
}
func rateErrorHandler(context echo.Context, err error) error {
	return context.JSON(http.StatusForbidden, nil)
}
func rateDenyHandler(context echo.Context, identifier string, err error) error {
	return context.JSON(http.StatusTooManyRequests, nil)
}

func main() {
	secretKey := utils.RandomBase64String(SECRET_LENGTH)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secretKey))))
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: RATE_TPS, Burst: BURST_TPS, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: identifierExtractor,
		ErrorHandler:        rateErrorHandler,
		DenyHandler:         rateDenyHandler,
	}
	e.Use(middleware.RateLimiterWithConfig(config))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogValuesFunc: logValuesFunc,
	}))
	e.Use(middleware.Recover())
	router.Router(e)
	fatal(start(":1323"))
}
