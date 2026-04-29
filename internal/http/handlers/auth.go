package handlers

import (
	"errors"
	"net/http"

	"chess-training/internal/http/middleware"
	"chess-training/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerReq struct {
	Username string  `json:"username" binding:"required,min=3,max=32"`
	Email    *string `json:"email" binding:"omitempty,email,max=254"`
	Password string  `json:"password" binding:"required,min=8,max=128"`
}

type loginReq struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required,min=3,max=254"`
	Password        string `json:"password" binding:"required,min=1,max=128"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid body"}})
		return
	}

	token, player, err := h.auth.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		status, code, message := registerErrorResponse(err)
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "player": newPlayerResponse(player)})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid body"}})
		return
	}

	token, player, err := h.auth.Login(c.Request.Context(), req.UsernameOrEmail, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "invalid credentials"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "player": newPlayerResponse(player)})
}

func Me(c *gin.Context) {
	response := meResponse{
		UserID:    contextString(c, middleware.CtxUserIDKey),
		Username:  contextString(c, middleware.CtxUsernameKey),
		RequestID: contextString(c, requestIDContextKey),
	}
	c.JSON(http.StatusOK, response)
}

type meResponse struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	RequestID string `json:"requestId,omitempty"`
}

const requestIDContextKey = "request_id"

func contextString(c *gin.Context, key string) string {
	raw, exists := c.Get(key)
	if !exists {
		return ""
	}
	value, ok := raw.(string)
	if !ok {
		return ""
	}
	return value
}

func registerErrorResponse(err error) (int, string, string) {
	switch {
	case errors.Is(err, service.ErrUsernameTaken):
		return http.StatusConflict, "USERNAME_TAKEN", "username already exists"
	case errors.Is(err, service.ErrEmailTaken):
		return http.StatusConflict, "EMAIL_TAKEN", "email already exists"
	case errors.Is(err, service.ErrInvalidUsername):
		return http.StatusBadRequest, "INVALID_USERNAME", "username must be 3-32 chars and contain only letters, numbers, or underscore"
	case errors.Is(err, service.ErrWeakPassword):
		return http.StatusBadRequest, "WEAK_PASSWORD", "password must be at least 8 characters"
	default:
		return http.StatusBadRequest, "REGISTER_FAILED", err.Error()
	}
}
