package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(routes *gin.Engine, PaymentController *controller.PaymentController, tokenRepo repository.TokenRepository) {
	authMiddleWare := middleware.AuthMiddleware(tokenRepo)

	paymentGroup := routes.Group("/payments")
	{
		// Protected routes (require valid authentication)
		paymentGroup.Use(authMiddleWare)
		{
			paymentGroup.POST("", PaymentController.CreatePayment)
			paymentGroup.PUT("/:id", PaymentController.UpdatePayment)
			paymentGroup.DELETE("/:id", PaymentController.DeletePayment)
			paymentGroup.GET("/:id", PaymentController.GetPaymentByID)
			paymentGroup.GET("", PaymentController.GetAllPayments)
		}

	}
	// WaafiPay webhook endpoint (no auth)
	routes.POST("/webhooks/waafi", PaymentController.HandleWaafiWebhook)

}
