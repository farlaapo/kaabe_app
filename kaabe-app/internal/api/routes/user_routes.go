package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user routes
func RegisterUserRoutes(router *gin.Engine, userController *controller.UserController, tokenRepo repository.TokenRepository) {
	// Apply middleware to protect certain routes
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	// User-related routes
	userGroup := router.Group("/users")
	{
		//Puplic routes
		userGroup.POST("", userController.RegisterUser)
		userGroup.POST("/authenticate", userController.AuthenticateUser)

		// Protected routes
		userGroup.Use(authMiddleware)
		{
			userGroup.GET("", userController.ListUsers)
			userGroup.GET("/:id", userController.GetUserByID)
			userGroup.PUT("/:id", userController.UpdateUser)
			userGroup.DELETE("/:id", userController.DeleteUser)
		}

	}
}
