package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRejectsTooShortPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(nil)
	router.POST("/auth/register", handler.Register)

	body := []byte(`{"username":"anouar","password":"1234567"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("BAD_REQUEST")) {
		t.Fatalf("expected BAD_REQUEST code in body, got %s", rec.Body.String())
	}
}

func TestLoginRejectsMissingUsernameOrEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(nil)
	router.POST("/auth/login", handler.Login)

	body := []byte(`{"password":"strong-enough"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("BAD_REQUEST")) {
		t.Fatalf("expected BAD_REQUEST code in body, got %s", rec.Body.String())
	}
}
