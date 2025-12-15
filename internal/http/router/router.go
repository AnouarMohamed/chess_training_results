package router

import (
	"chess-training/internal/http/handlers"
	"chess-training/internal/http/middleware"
	"chess-training/internal/service"

	"github.com/gin-gonic/gin"
)

func New(authSvc *service.AuthService, jwtSecret string) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", handlers.Healthz)

	authHandler := handlers.NewAuthHandler(authSvc)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.RequireAuth(jwtSecret))
	protected.GET("/auth/me", handlers.Me)

	return r
}
