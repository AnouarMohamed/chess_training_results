package middleware

import (
	"errors"
	"net/http"
	"strings"

	"chess-training/internal/util"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "user_id"
const CtxUsernameKey = "username"

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := bearerTokenFromAuthorizationHeader(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing bearer token"}})
			return
		}
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

func bearerTokenFromAuthorizationHeader(header string) (string, error) {
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return "", errors.New("authorization header must use bearer token")
	}
	if strings.TrimSpace(fields[1]) == "" {
		return "", errors.New("bearer token is empty")
	}
	return fields[1], nil
}
