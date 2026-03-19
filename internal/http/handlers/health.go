package handlers

import (
	"chess-training/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	requestID, _ := c.Get(middleware.CtxRequestIDKey)
	c.JSON(200, gin.H{
		"ok":        true,
		"requestId": requestID,
	})
}
