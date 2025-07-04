package controller

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
 
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
)

// withdrawalController struct that defines the withdrawal controller with its service
type WithdrawalController struct {
	WithdrawalService service.WithdrawalService
}

func NewWithdrawalController(withdrawalService service.WithdrawalService) *WithdrawalController {
	return &WithdrawalController{WithdrawalService: withdrawalService}
}

func (wc *WithdrawalController) CreateWithdrawal(ctx * gin.Context) {
	 var Withdrawal model.Withdrawal

	// Bind JSON to the createdWithdrawal struct, which now includes VideoUrl as []string
	if err := ctx.ShouldBindJSON(&Withdrawal); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the service with the bound struct's videoUrl slice and WithdrawalD
	createdWithdrawal, err := wc.WithdrawalService.CreateWithdrawal(Withdrawal.InfluencerID, Withdrawal.Amount, Withdrawal.Status, Withdrawal.RequestedAt, Withdrawal.ProcessedAt )

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, createdWithdrawal)
}

func (wc *WithdrawalController) UpdateWithdrawal(ctx * gin.Context) {
	
	var Withdrawal model.Withdrawal

	WithdrawalIdParam := ctx.Param("id")
	WithdrawalID, err := uuid.FromString(WithdrawalIdParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid Withdrawal ID"})
		return
	}

	// Bind JSON to the Withdrawal struct, which now
	if err := ctx.ShouldBindJSON(&Withdrawal); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	Withdrawal.ID = WithdrawalID

	// Call the service with the bound and Withdrawal
	if err := wc.WithdrawalService.UpdateWithdrawal(&Withdrawal); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Withdrawal updated successfully"})
 }

func (wc *WithdrawalController) DeleteWithdrawal(ctx * gin.Context) {
		WithdrawalParam := ctx.Param("id")
	 WithdrawalId, err := uuid.FromString(WithdrawalParam)
	 if err != nil {
	 	ctx.JSON(400, gin.H{"error": "invalid Withdrawal ID"})
	}

	if err := wc.WithdrawalService.DeleteWithdrawal(WithdrawalId); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Withdrawal deleted successfully"})
	}

func (wc *WithdrawalController) GetWithdrawalByID(ctx * gin.Context) {
	WithdrawalParam := ctx.Param("id")
	log.Printf("Lesson ID from URL: %s", WithdrawalParam)
	WithdrawalID, err := uuid.FromString(WithdrawalParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid lesson ID format"})
		return
	}

	Withdrawal, err := wc.WithdrawalService.GetWithdrawalByID(WithdrawalID)
	if err != nil {
		if err.Error() == "Withdrawal not found" {
			ctx.JSON(404, gin.H{"error": err.Error()})
		} else {
			log.Printf("Unexpected error: %v", err)
			ctx.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}

	ctx.JSON(200, Withdrawal)
}
 
func (wc *WithdrawalController) GetAllWithdrawal(ctx * gin.Context) {
	// Call service to get Withdrawal
	Withdrawal, err := wc.WithdrawalService.GetAllWithdrawal()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// respond success
	ctx.JSON(200, Withdrawal)
}
