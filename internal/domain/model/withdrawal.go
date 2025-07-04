package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Withdrawal struct {
	ID           uuid.UUID `json:"id"`
	InfluencerID uuid.UUID `json:"influencer_id"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"` // e.g., "pending", "approved", "rejected"
	RequestedAt  time.Time `json:"requested_at"`
	ProcessedAt  time.Time  `json:"processed_at,omitempty"` // Nullable, if not yet processed
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
