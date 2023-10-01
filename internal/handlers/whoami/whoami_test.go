package whoami

import (
	"encoding/json"
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

func TestWhoAmINoUserID(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secretKey"))))
	req := httptest.NewRequest(http.MethodGet, "/ping", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sessGet = func(name string, c echo.Context) (*sessions.Session, error) {
		return sessions.NewSession(MockStore{}, "mockStore"), nil
	}
	WhoAmI(c)
	resp := map[string]string{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["message"] != "no session" {
		t.Fatalf("expected no session: %s", resp["message"])
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request: %d", rec.Code)
	}
}

func TestWhoAmIUserID(t *testing.T) {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secretKey"))))
	req := httptest.NewRequest(http.MethodGet, "/ping", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	sessGet = func(name string, c echo.Context) (*sessions.Session, error) {
		sess := sessions.NewSession(MockStore{}, "mockStore")
		sess.Values["user"] = "techytechster"
		return sess, nil
	}
	WhoAmI(c)
	resp := map[string]string{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected ok request: %d", rec.Code)
	}
	if resp["userId"] != "techytechster" {
		t.Fatalf("expected techytechster as userid: %s", resp["userId"])
	}
}
