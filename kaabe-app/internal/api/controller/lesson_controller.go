package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// LessonController struct that defines the lesson controller with its service
type LessonController struct {
	LessonService service.LessonService
}

// NewLessonController creates a new LessonController instance
func NewLessonController(lessonService service.LessonService) *LessonController {
	return &LessonController{LessonService: lessonService}
}

// CreateLesson handles the creation of a new lesson
func (l *LessonController) CreateLesson(ctx *gin.Context) {
	var lesson model.Lesson

	// Bind JSON to the lesson struct, which now includes VideoUrl as []string
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service with the bound struct's videoUrl slice and lessonID
	createdLesson, err := l.LessonService.CreateLesson(
		lesson.CourseID,
		lesson.Title,
		lesson.VideoURL,
		lesson.Order,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdLesson)

}

// UpdateLesson handles the update of an existing lesson
func (l *LessonController) UpdateLesson(ctx *gin.Context) {
	var lesson model.Lesson

	lessonIdParam := ctx.Param("id")
	lessonID, err := uuid.FromString(lessonIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	// Bind JSON to the lesson struct, which now includes videoUrl as []string
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson.ID = lessonID

	// Call the service with the bound struct's videoUrl slice and lessonID
	if err := l.LessonService.UpdateLesson(&lesson); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "lesson updated successfully"})

}

// DeleteLesson handles the delete of an existing lesson
func (l *LessonController) DeleteLesson(ctx *gin.Context) {

	lessonIdParam := ctx.Param("id")
	lessonId, err := uuid.FromString(lessonIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Lesson ID"})
	}

	if err := l.LessonService.DeleteLesson(lessonId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "lesson deleted successfully"})

}

// get lesson by id
func (l *LessonController) GetLessonByID(ctx *gin.Context) {
	lessonParam := ctx.Param("id")
	log.Printf("Lesson ID from URL: %s", lessonParam)

	lessonID, err := uuid.FromString(lessonParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson ID format"})
		return
	}

	lesson, err := l.LessonService.GetLessonByID(lessonID)
	if err != nil {
		if err.Error() == "lesson not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unexpected error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(http.StatusOK, lesson)
}


// get all lesson
func (l *LessonController) GetAllLessons(ctx *gin.Context) {
	// Call service to get lesson
	lesson, err := l.LessonService.GetAllLessons()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// respond success
	ctx.JSON(http.StatusOK, lesson)
}
