package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email" binding:"required,email"`
	Password         string     `json:"password" binding:"required,min=6"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Role             string     `json:"role"`
	WalletID         *string    `json:"wallet_id,omitempty"` // omit if nil
	ResetToken       *uuid.UUID `json:"reset_token,omitempty"`        // omit if not used
	ResetTokenExpiry *time.Time `json:"reset_token_expiry,omitempty"` // omit if not used
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
