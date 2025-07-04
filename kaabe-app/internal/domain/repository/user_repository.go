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
}