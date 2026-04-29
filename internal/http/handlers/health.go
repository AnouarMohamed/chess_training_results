package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Healthz(c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	c.JSON(200, gin.H{"ok": true, "env": env})
}
