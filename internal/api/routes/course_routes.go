package routes

import (
	"kaabe-app/internal/api/controller"
	"kaabe-app/internal/api/middleware"
	"kaabe-app/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

func RegisterCoursesRoutes(routes *gin.Engine, courseController *controller.CourseController, tokenRepo repository.TokenRepository) {
	authMiddleware := middleware.AuthMiddleware(tokenRepo)

	courseGroup := routes.Group("/courses")
	{

		// Protected routes (require valid authentication)
		courseGroup.Use(authMiddleware)
		{
			courseGroup.POST("", courseController.CreateCourse)
			courseGroup.PUT("/:id", courseController.UpdateCourse)
			courseGroup.DELETE("/:id", courseController.DeleteCourse)
			courseGroup.GET("/:id", courseController.GetCourseByID)
			courseGroup.GET("", courseController.GetAllCourses)
		}

	}
}
