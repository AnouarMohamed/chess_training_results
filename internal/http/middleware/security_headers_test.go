package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSecurityHeadersAreSet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertHeader := func(key, want string) {
		if got := rec.Header().Get(key); got != want {
			t.Fatalf("expected %s=%q, got %q", key, want, got)
		}
	}

	assertHeader("X-Content-Type-Options", "nosniff")
	assertHeader("X-Frame-Options", "DENY")
	assertHeader("Referrer-Policy", "no-referrer")
	assertHeader("X-XSS-Protection", "0")
}
