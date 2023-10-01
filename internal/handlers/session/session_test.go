package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type MockStore struct {
	sessions.Store
}

func (m MockStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

func TestAcquireSession(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secretKey"))))
	req := httptest.NewRequest(http.MethodGet, "/session", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sessGet = func(name string, c echo.Context) (*sessions.Session, error) {
		return sessions.NewSession(MockStore{}, "mockStore"), nil
	}
	AcquireSession(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected statusTeapot: %d", rec.Code)
	}
}
