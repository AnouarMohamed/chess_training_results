package handlers

import (
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
	Username string  `json:"username" binding:"required"`
	Email    *string `json:"email"`
	Password string  `json:"password" binding:"required"`
}

type loginReq struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid body"}})
		return
	}

	token, player, err := h.auth.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		code := "REGISTER_FAILED"
		status := http.StatusBadRequest
		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
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
