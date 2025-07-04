package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// ratingController struct that defines the rating controller with its service
type RatingController struct {
	RatingService service.RatingService
}

// NewRatingController creates a new ratingController instance
func NewRatingController(ratingService service.RatingService) *RatingController {
	return &RatingController{RatingService: ratingService}
}

// CreateRating handles the creation of a new rating
func (r *RatingController) CreateRating(ctx *gin.Context) {
	var rating model.Rating

	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdRating, err := r.RatingService.CreateRating(
		rating.UserID,
		rating.CourseID,
		rating.Score,
		rating.Comment,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdRating)
}

func (r *RatingController) UpdateRating(ctx *gin.Context) {
	var rating model.Rating

	ratingIdParam := ctx.Param("id")
	ratingID, err := uuid.FromString(ratingIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	//
	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating.ID = ratingID

	if err := r.RatingService.UpdateRating(&rating); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "rating updated successfully"})
}

func (r *RatingController) DeleteRating(ctx *gin.Context) {
	ratingIdParam := ctx.Param("id")
	ratingId, err := uuid.FromString(ratingIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating ID"})
	}

	if err := r.RatingService.DeleteRating(ratingId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "rating deleted successfully"})

}

func (r *RatingController) GetAllRating(ctx *gin.Context) {
	// Call service to get rating
	rating, err := r.RatingService.GetAllRatings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// respond success
	ctx.JSON(http.StatusOK, rating)
}

// GetRatingByID godoc
func (r *RatingController) GetRatingByID(c *gin.Context) {
	ratingParam := c.Param("id")
	log.Printf("rating ID : %s", ratingParam)

	ratingID, err := uuid.FromString(ratingParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rating ID format"})
		return
	}

	rating, err := r.RatingService.GetRatingByID(ratingID)
	if err != nil {
		if err.Error() == "rating not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, rating)
}
