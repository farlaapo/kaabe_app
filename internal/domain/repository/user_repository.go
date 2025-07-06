package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(userID uuid.UUID) error
	Get(userID uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	List() ([]*model.User, error)

	// Password reset
	SetResetToken(email string, token uuid.UUID, expiry string) error
	FindByResetToken(token uuid.UUID) (*model.User, error)
	UpdatePassword(userID uuid.UUID, hashedPassword string) error
	ClearResetToken(userID uuid.UUID) error
}
