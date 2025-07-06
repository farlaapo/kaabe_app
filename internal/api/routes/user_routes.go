package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(router *gin.Engine, userController *controller.UserController, tokenRepo repository.TokenRepository) {
	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	// Group all /users endpoints
	userGroup := router.Group("/users")
	{
		// 🚪 Public Routes
		userGroup.POST("", userController.RegisterUser)
		userGroup.POST("/authenticate", userController.AuthenticateUser)
		userGroup.POST("/forgot-password", userController.ForgotPassword)
		userGroup.POST("/reset-password", userController.ResetPassword)

		// 🔒 Protected Routes (Require Auth)
		userGroup.Use(authMiddleware)
		{
			userGroup.GET("", userController.ListUsers)
			userGroup.GET("/:id", userController.GetUserByID)
			userGroup.PUT("/:id", userController.UpdateUser)
			userGroup.DELETE("/:id", userController.DeleteUser)
		}
	}
}
