package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

type RatingRepository interface {
	Create(rating *model.Rating) error
	Update(rating *model.Rating) error
	Delete(ratingID uuid.UUID) error
	GetByID(ratingID uuid.UUID) (*model.Rating, error)
	Getall() ([]*model.Rating, error)
}
