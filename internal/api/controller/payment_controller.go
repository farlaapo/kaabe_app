package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type PaymentController struct {
	paymentService service.PaymentService
	//walletService  service.WalletService
}

// NewPaymentController creates a new PaymentController instance
func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService} //walletService:  walletService,

}

// CreatePayment handles the creation of a new payment
func (p *PaymentController) CreatePayment(ctx *gin.Context) {
	var payment model.Payment

	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPayment, err := p.paymentService.CreatePayment(
		payment.ExternalRef,
		payment.UserID,
		payment.SubscriptionID,
		payment.Amount,
		payment.Status,
		payment.ProcessedAt,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdPayment)
}

// UpdatePayment handles updating an existing payment
func (p *PaymentController) UpdatePayment(ctx *gin.Context) {
	var payment model.Payment

	paymentIDParam := ctx.Param("id")
	paymentID, err := uuid.FromString(paymentIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment.ID = paymentID
	payment.UpdatedAt = time.Now()

	if err := p.paymentService.UpdatePayment(&payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment updated successfully"})
}

// DeletePayment handles the soft deletion of a payment
func (p *PaymentController) DeletePayment(ctx *gin.Context) {
	paymentIDParam := ctx.Param("id")
	paymentID, err := uuid.FromString(paymentIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := p.paymentService.DeletePayment(paymentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment deleted successfully"})
}

// GetPaymentByID retrieves a payment by its ID
func (p *PaymentController) GetPaymentByID(ctx *gin.Context) {
	paymentIDParam := ctx.Param("id")
	log.Printf("Payment ID from URL: %s", paymentIDParam)

	paymentID, err := uuid.FromString(paymentIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID format"})
		return
	}

	payment, err := p.paymentService.GetPaymentByID(paymentID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payment)
}

// GetAllPayments retrieves all non-deleted payments
func (p *PaymentController) GetAllPayments(ctx *gin.Context) {
	payments, err := p.paymentService.GetAllPayments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payments)
}

// WaafiPay Webhook
func (p *PaymentController) HandleWaafiWebhook(c *gin.Context) {
	var payload struct {
		ExternalRef    string  `json:"external_ref"`
		UserID         string  `json:"user_id"`
		SubscriptionID string  `json:"subscription_id"`
		Amount         float64 `json:"amount"`
		Status         string  `json:"status"`
		ProcessedAt    string  `json:"processed_at"`
		Signature      string  `json:"signature"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Signature Verification Implementation (example HMAC)
	expectedSecret := "your-waafi-secret"
	computedSignature := ComputeHMAC(payload.ExternalRef+payload.UserID+payload.SubscriptionID+fmt.Sprintf("%f", payload.Amount)+payload.Status+payload.ProcessedAt, expectedSecret)
	if payload.Signature != computedSignature {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// ✅ Duplicate Check
	existing, err := p.paymentService.GetPaymentByExternalRef(payload.ExternalRef)
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate external_ref"})
		return
	}

	// Retry-safe parse
	userID, err := uuid.FromString(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}
	subID, err := uuid.FromString(payload.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription_id"})
		return
	}
	processedAt, err := time.Parse(time.RFC3339, payload.ProcessedAt)
	if err != nil {
		processedAt = time.Now()
	}

	payment, err := p.paymentService.CreatePayment(payload.ExternalRef, userID, subID, payload.Amount, payload.Status, processedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// ComputeHMAC generates an HMAC SHA256 hash
func ComputeHMAC(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// // WaafiPay Webhook
// // Payments: Wallet integrations (initial target: WaafiPay); webhook‑driven events
// func (p *PaymentController) HandleWaafiWebhook(c *gin.Context) {
// 	var payload struct {
// 		ExternalRef    string  `json:"external_ref"`
// 		UserID         string  `json:"user_id"`
// 		SubscriptionID string  `json:"subscription_id"`
// 		Amount         float64 `json:"amount"`
// 		Status         string  `json:"status"`
// 		ProcessedAt    string  `json:"processed_at"`
// 		Signature      string  `json:"signature"`
// 	}
// 	if err := c.ShouldBindJSON(&payload); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// ✅ Signature Verification Placeholder
// 	if payload.Signature != "expected-signature-value" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
// 		return
// 	}

// 	// ✅ Duplicate Check
// 	existing, err := p.paymentService.GetPaymentByExternalRef(payload.ExternalRef)
// 	if err == nil && existing != nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate external_ref"})
// 		return
// 	}

// 	// Retry-safe parse
// 	userID, err := uuid.FromString(payload.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
// 		return
// 	}
// 	subID, err := uuid.FromString(payload.SubscriptionID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription_id"})
// 		return
// 	}
// 	processedAt, err := time.Parse(time.RFC3339, payload.ProcessedAt)
// 	if err != nil {
// 		processedAt = time.Now()
// 	}

// 	payment, err := p.paymentService.CreatePayment(payload.ExternalRef, userID, subID, payload.Amount, payload.Status, processedAt)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, payment)
// }
