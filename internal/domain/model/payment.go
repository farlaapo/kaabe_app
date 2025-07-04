package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Payment struct {
	ID             uuid.UUID `json:"id"`
	ExternalRef    string    `json:"external_ref"`
	UserID         uuid.UUID `json:"user_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`       // e.g., "pending", "completed", "failed"
	ProcessedAt    time.Time `json:"processed_at"` // time the payment was processed
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
