package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type SubscriptionController struct {
	SubscriptionService service.SubscriptionService
}

func NewSubscriptionController(SubscriptionService service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{SubscriptionService: SubscriptionService}
}

func (sc *SubscriptionController) CreateSubscription(ctx *gin.Context) {
	var Subscription model.Subscription

	// Bind JSON to the  struct,
	if err := ctx.ShouldBindJSON(&Subscription); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the service with the
	createdSubscription, err := sc.SubscriptionService.CreateSubscription(
		Subscription.UserID,
		Subscription.CourseID,
		Subscription.StartedAt,
		Subscription.ExpiresAt,
		Subscription.Status,
	)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, createdSubscription)
}

func (sc *SubscriptionController) UpdateSubscription(ctx *gin.Context) {
	var Subscription model.Subscription

	SubscriptionIdParam := ctx.Param("id")
	SubscriptionID, err := uuid.FromString(SubscriptionIdParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Subscription ID"})
		return
	}

	// Bind JSON to the  struct, which now includes
	if err := ctx.ShouldBindJSON(&Subscription); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	Subscription.ID = SubscriptionID

	// Call the service with the bound struct's
	if err := sc.SubscriptionService.UpdateSubscription(&Subscription); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Subscription updated successfully"})
}

func (sc *SubscriptionController) DeleteSubscription(ctx *gin.Context) {
	SubscriptionParam := ctx.Param("id")
	SubscriptionId, err := uuid.FromString(SubscriptionParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid Subscription ID"})
	}

	if err := sc.SubscriptionService.DeleteSubscription(SubscriptionId); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Subscription deleted successfully"})
}

func (sc *SubscriptionController) GetSubscriptionByID(ctx *gin.Context) {
	SubscriptionParam := ctx.Param("id")
	log.Printf("Lesson ID from URL: %s", SubscriptionParam)
	subscriptionID, err := uuid.FromString(SubscriptionParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid lesson ID format"})
		return
	}

	Subscription, err := sc.SubscriptionService.GetSubscriptionByID(subscriptionID)
	if err != nil {
		if err.Error() == "Subscription not found" {
			ctx.JSON(404, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unexpected error: %v", err)
			ctx.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(200, Subscription)
}

func (sc *SubscriptionController) GetAllSubscription(ctx *gin.Context) {
	// Call service to get Subscription
	Subscription, err := sc.SubscriptionService.GetAllSubscription()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// respond success
	ctx.JSON(http.StatusOK, Subscription)
}
