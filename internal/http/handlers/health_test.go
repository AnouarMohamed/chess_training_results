package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthzDefaultsToDevEnv(t *testing.T) {
	t.Setenv("APP_ENV", "")
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/healthz", nil)

	Healthz(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if payload["ok"] != true {
		t.Fatalf("expected ok=true, got %v", payload["ok"])
	}
	if payload["env"] != "dev" {
		t.Fatalf("expected env=dev, got %v", payload["env"])
	}
}

func TestHealthzUsesConfiguredEnv(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/healthz", nil)

	Healthz(ctx)

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if payload["env"] != "staging" {
		t.Fatalf("expected env=staging, got %v", payload["env"])
	}
}
