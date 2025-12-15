package middleware

import (
	"net/http"
	"strings"

	"chess-training/internal/util"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"
const CtxUsernameKey = "username"

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing bearer token"}})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := util.ParseJWT(jwtSecret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "invalid token"}})
			return
		}
		c.Set(CtxUserIDKey, claims.Sub)
		c.Set(CtxUsernameKey, claims.Username)
		c.Next()
	}
}
