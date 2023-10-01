package router

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRouter(t *testing.T) {
	Router(echo.New())
}
