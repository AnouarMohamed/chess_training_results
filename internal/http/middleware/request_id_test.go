package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGeneratesWhenHeaderMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/ping", func(c *gin.Context) {
		value, _ := c.Get(CtxRequestIDKey)
		c.JSON(http.StatusOK, gin.H{"requestId": value})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	headerValue := rec.Header().Get(HeaderRequestID)
	if headerValue == "" {
		t.Fatalf("expected response %s header to be set", HeaderRequestID)
	}
	if !strings.Contains(rec.Body.String(), headerValue) {
		t.Fatalf("expected body to contain generated request id %q, got %s", headerValue, rec.Body.String())
	}
}

func TestRequestIDUsesIncomingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/ping", func(c *gin.Context) {
		value, _ := c.Get(CtxRequestIDKey)
		c.JSON(http.StatusOK, gin.H{"requestId": value})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(HeaderRequestID, "custom-req-id")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if got := rec.Header().Get(HeaderRequestID); got != "custom-req-id" {
		t.Fatalf("expected %s header %q, got %q", HeaderRequestID, "custom-req-id", got)
	}
	if !strings.Contains(rec.Body.String(), "custom-req-id") {
		t.Fatalf("expected body to contain provided request id, got %s", rec.Body.String())
	}
}
