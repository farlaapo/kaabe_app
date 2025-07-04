package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func RegisterLessonRoutes(routes *gin.Engine, lessonController *controller.LessonController, tokenRepo repository.TokenRepository) {
	authMiddleWare := middleware.AuthMiddleware(tokenRepo)

	courGroup := routes.Group("/lessons")
	{
		// Protected routes (require valid authentication)
		courGroup.Use(authMiddleWare)
		{
			courGroup.POST("", lessonController.CreateLesson)
			courGroup.PUT("/:id", lessonController.UpdateLesson)
			courGroup.DELETE("/:id", lessonController.DeleteLesson)
			courGroup.GET("/:id", lessonController.GetLessonByID)
			courGroup.GET("", lessonController.GetAllLessons)
		}
	}

}
