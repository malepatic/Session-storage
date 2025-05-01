package api

import (
	"session-app/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService service.AuthService) {
	handler := NewHandler(authService)

	// Public routes
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.GET("/health", handler.HealthCheck)

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware(authService))
	{
		protected.POST("/logout", handler.Logout)
		protected.GET("/profile", handler.GetProfile)
	}
}
