package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chess-training/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func TestMeIncludesIdentityAndRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	ctx.Set(middleware.CtxUserIDKey, "user-99")
	ctx.Set(middleware.CtxUsernameKey, "anouar")
	ctx.Set(requestIDContextKey, "req-abc")

	Me(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	if payload["userId"] != "user-99" {
		t.Fatalf("expected userId user-99, got %v", payload["userId"])
	}
	if payload["username"] != "anouar" {
		t.Fatalf("expected username anouar, got %v", payload["username"])
	}
	if payload["requestId"] != "req-abc" {
		t.Fatalf("expected requestId req-abc, got %v", payload["requestId"])
	}
}

func TestMeOmitsRequestIDWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	ctx.Set(middleware.CtxUserIDKey, "user-77")
	ctx.Set(middleware.CtxUsernameKey, "reader")

	Me(ctx)

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	if _, exists := payload["requestId"]; exists {
		t.Fatalf("expected requestId to be omitted when missing")
	}
}
