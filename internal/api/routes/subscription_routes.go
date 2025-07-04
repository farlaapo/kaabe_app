package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func RegisterSubscriptionRoutes(routes *gin.Engine, SubscriptionController *controller.SubscriptionController, tokenRepo repository.TokenRepository) {
	authMiddleWare := middleware.AuthMiddleware(tokenRepo)

	courGroup := routes.Group("/Subscription")
	{
		// Protected routes (require valid authentication)
		courGroup.Use(authMiddleWare)
		{
			courGroup.POST("", SubscriptionController.CreateSubscription)
			courGroup.PUT("/:id", SubscriptionController.UpdateSubscription)
			courGroup.DELETE("/:id", SubscriptionController.DeleteSubscription)
			courGroup.GET("/:id", SubscriptionController.GetSubscriptionByID)
			courGroup.GET("", SubscriptionController.GetAllSubscription)
		}
	}

}