package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chess-training/internal/util"
	"github.com/gin-gonic/gin"
)

func TestRequireAuthRejectsMissingBearerHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", RequireAuth("secret"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAuthRejectsInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", RequireAuth("secret"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestRequireAuthAllowsValidTokenAndSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "secret"
	token, err := util.SignJWT(secret, "user-123", "anouar", 30)
	if err != nil {
		t.Fatalf("SignJWT returned error: %v", err)
	}

	router := gin.New()
	router.GET("/protected", RequireAuth(secret), func(c *gin.Context) {
		userID, _ := c.Get(CtxUserIDKey)
		username, _ := c.Get(CtxUsernameKey)
		c.JSON(http.StatusOK, gin.H{
			"userId":   userID,
			"username": username,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "bearer   "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "user-123") || !strings.Contains(body, "anouar") {
		t.Fatalf("expected response to include user context, got %s", body)
	}
}
