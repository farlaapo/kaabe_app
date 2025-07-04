package service

import (
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type RatingService interface {
	CreateRating(UserID, CourseID uuid.UUID, Score int, Comment string) (*model.Rating, error)
	UpdateRating(rating model.Rating) error
	DeleteRating(RatingID uuid.UUID) error
	GetRatingByID(ratingID uuid.UUID) (*model.Rating, error)
	GetAllRatings() ([]*model.Rating, error)
}

// RatingServiceImpl struct implementing ratingService
type RatingServiceImpl struct {
	repo      repository.RatingRepository
	tokenRepo repository.TokenRepository
}

// CreateRating implements RatingService.
func (r *RatingServiceImpl) CreateRating(UserID uuid.UUID, CourseID uuid.UUID, Score int, Comment string) (*model.Rating, error) {
	// Generate a new UUID for the rating ID

	neoRating, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Create a new rating instance
	amRating := &model.Rating{
		ID:        neoRating,
		UserID:    UserID,
		CourseID:  CourseID,
		Score:     Score,
		Comment:   Comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Log the new rating creation attempt
	log.Printf("Creating rating: %+v", amRating)

	// Save the new rating to the repository
	err = r.repo.Create(amRating)
	if err != nil {
		return nil, fmt.Errorf("failed to create rating: %v", err)
	}

	return amRating, nil
}

// DeleteRating implements RatingService.
func (r *RatingServiceImpl) DeleteRating(RatingID uuid.UUID) error {
	// Find the rating to delete
	_, err := r.repo.GetByID(RatingID)
	if err != nil {
		return fmt.Errorf("could not find rating with ID %s: %v", RatingID, err)
	}

	if err := r.repo.Delete(RatingID); err != nil {
		return fmt.Errorf("failed to delete rating with ID %s: %v", RatingID, err)
	}

	log.Printf("Successfully deleted rating with ID %s", RatingID)
	return nil
}

// GetAllRatings implements RatingService.
func (r *RatingServiceImpl) GetAllRatings() ([]*model.Rating, error) {
	// Retrieve all ratings from the repository
	rating, err := r.repo.Getall()
	if err != nil {
		return nil, fmt.Errorf("failed to get all rating: %v", err)
	}
	return rating, nil

}

// GetRatingByID implements RatingService.
func (r *RatingServiceImpl) GetRatingByID(ratingID uuid.UUID) (*model.Rating, error) {
	// Retrieve the rating from the repository
	rating, err := r.repo.GetByID(ratingID)
	if err != nil {
		return nil, err
	}
	return rating, nil
}

// UpdateRating implements RatingService.
func (r *RatingServiceImpl) UpdateRating(rating model.Rating) error {
	// Find the rating to update
	_, err := r.repo.GetByID(rating.ID)
	if err != nil {
		return fmt.Errorf("could not find rating with ID %s", rating.ID)
	}

	if err := r.repo.Update(&rating); err != nil {
		return fmt.Errorf("failed to update lesson with ID %s: %v", rating.ID, err)
	}

	return nil
}

func NewRatingService(ratingRepo repository.RatingRepository, tokenRepo repository.TokenRepository) RatingService {
	return &RatingServiceImpl{
		repo:      ratingRepo,
		tokenRepo: tokenRepo,
	}
}
