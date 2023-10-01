package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TestMain(t *testing.T) {
	fatal = func(i ...interface{}) {}
	start = func(address string) error {
		return nil
	}
	main()
}

func TestLogValuesFunc(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Jon Snow","email":"jon@labstack.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	logValuesFunc(e.NewContext(req, rec), middleware.RequestLoggerValues{})
}

func TestRateLimitHelpers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/session", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	identifierExtractor(c)
	rateErrorHandler(c, nil)
	rateDenyHandler(c, "uniqueIdentifier", nil)
}
