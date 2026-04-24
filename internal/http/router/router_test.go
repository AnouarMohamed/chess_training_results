package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chess-training/internal/util"
	"github.com/gin-gonic/gin"
)

func TestHealthzRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(nil, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"ok":true`) {
		t.Fatalf("expected health response to contain ok=true, got %s", rec.Body.String())
	}
}

func TestAuthMeRequiresToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(nil, "test-secret")

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestAuthMeWithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	r := New(nil, secret)

	token, err := util.SignJWT(secret, "user-42", "anouar", 30)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "user-42") || !strings.Contains(body, "anouar") {
		t.Fatalf("expected response to contain token identity, got %s", body)
	}
}
