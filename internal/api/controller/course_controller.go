package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CourseController struct that defines the course controller with its service
type CourseController struct {
	courseService service.CourseService
}

// NewCourseController creates a new CourseController instance
func NewCourseController(courseService service.CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

// CreateCourse handles the creation of a new course
func (c *CourseController) CreateCourse(ctx *gin.Context) {
	var course model.Course

	// Bind JSON to the course struct, which now includes CoverImageURL as []string
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the request  set by the AuthMiddleware
	_, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is required"})
		return
	}

	// Call the service with the bound struct's CoverImageURL slice and instructorID
	createdCourse, err := c.courseService.CreateCourse(
		course.InfluencerID,
		course.Title,
		course.Description,
		course.Price,
		course.CoverImageURL,
		course.Status,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdCourse)
}

// UpdateCourse handles the update of an existing course
func (c *CourseController) UpdateCourse(ctx *gin.Context) {
	var course model.Course

	// Parse and validate course ID from URL
	courseIdParam := ctx.Param("id")
	courseID, err := uuid.FromString(courseIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Bind JSON to the course struct, which now includes CoverImageURL as []string
	if err := ctx.ShouldBindJSON(&course); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.ID = courseID

	// Call the service with the bound struct's CoverImageURL slice and instructorID
	if err := c.courseService.UpdateCourse(&course); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DeleteCourse handles the delete of an existing course
func (c *CourseController) DeleteCourse(ctx *gin.Context) {

	// Parse and validate course ID from URL
	courseIdParam := ctx.Param("id")
	courseID, err := uuid.FromString(courseIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Call service to delete course
	if err := c.courseService.DeleteCourse(courseID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// get course by id
func (c *CourseController) GetCourseByID(ctx *gin.Context) {

	// Parse and validate course ID from URL
	courseIdParam := ctx.Param("id")
	courseID, err := uuid.FromString(courseIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Call service to get course
	course, err := c.courseService.GetCourseByID(courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, course)

}

// getall courses
func (c *CourseController) GetAllCourses(ctx *gin.Context) {
	// Call service to get course
	courses, err := c.courseService.GetAllCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// respond success
	ctx.JSON(http.StatusOK, courses)

}
