package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)



func RegisterWithdrawalRoutes(routes *gin.Engine, WithdrawalController *controller.WithdrawalController, tokenRepo repository.TokenRepository) {
	authMiddleWare := middleware.AuthMiddleware(tokenRepo)

	courGroup := routes.Group("/Withdrawal")
	{
		// Protected routes (require valid authentication)
		courGroup.Use(authMiddleWare)
		{
			courGroup.POST("", WithdrawalController.CreateWithdrawal)
			courGroup.PUT("/:id", WithdrawalController.UpdateWithdrawal)
			courGroup.DELETE("/:id", WithdrawalController.DeleteWithdrawal)
			courGroup.GET("/:id", WithdrawalController.GetWithdrawalByID)
			courGroup.GET("", WithdrawalController.GetAllWithdrawal)
		}
	}

}