package ping

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHeartbeat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	Ping(c)
	if rec.Code != http.StatusTeapot {
		t.Fatalf("expected statusTeapot: %d", rec.Code)
	}
	resp := map[string]string{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["message"] != "pong" {
		t.Fatalf("expected message = pong as response recieved: %s", rec.Body.String())
	}
}
