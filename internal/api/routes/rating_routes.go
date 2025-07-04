package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func RegisterRatingRoutes(routes *gin.Engine, ratingcontroller *controller.RatingController, tokenRepo repository.TokenRepository) {
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	ratingGroup := routes.Group("/rating")
	{

		// Protected routes (require valid authentication)
		ratingGroup.Use(authMiddleware)
		{
			ratingGroup.POST("/", ratingcontroller.CreateRating)
			ratingGroup.PUT("/:id", ratingcontroller.UpdateRating)
			ratingGroup.DELETE("/:id", ratingcontroller.DeleteRating)
			ratingGroup.GET("/:id", ratingcontroller.GetRatingByID)
			ratingGroup.GET("", ratingcontroller.GetAllRating)
		}
	}
}
