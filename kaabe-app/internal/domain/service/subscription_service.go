package service

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"time"
 "log"
 "fmt"
	"github.com/gofrs/uuid"
)

type SubscriptionService interface{
	CreateSubscription(userID uuid.UUID, courseID uuid.UUID, StartedAt time.Time , ExpiresAt time.Time, status string ) (*model.Subscription, error)
	UpdateSubscription(Subscription  *model.Subscription) error
  GetSubscriptionByID(SubscriptionID uuid.UUID) (*model.Subscription, error)
	GetAllSubscription()([]*model.Subscription, error)
	DeleteSubscription(SubscriptionID uuid.UUID) error
} 

// lessonServiceImpl struct implementing lessonService
type SubscriptionServiceImpl struct {
	repo      repository.SubscriptionRepository
	tokenRepo repository.TokenRepository
}

func NewSubscriptionService(SubscriptionRep repository.SubscriptionRepository, tokenRepo repository.TokenRepository) SubscriptionService {
  return &SubscriptionServiceImpl{
		repo:      SubscriptionRep,
		tokenRepo: tokenRepo,
	}
}

func (s *SubscriptionServiceImpl) CreateSubscription(userID uuid.UUID, courseID uuid.UUID, StartedAt time.Time , ExpiresAt time.Time, status string ) (*model.Subscription, error) {
	// Generate a new UUID for the lesson ID

	neoLesson, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Create a new lesson instance
	amsubs := &model.Subscription{
		ID:        neoLesson,
		UserID: userID,
		CourseID:  courseID,
		StartedAt: StartedAt,
		ExpiresAt:   ExpiresAt,
		Status: status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Log the new lesson creation attempt
	log.Printf("Creating lesson: %+v", amsubs)

	// Save the new lesson to the repository
	err = s.repo.Create(amsubs)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson: %v", err)
	}

	return amsubs, nil
}

func (s *SubscriptionServiceImpl) DeleteSubscription(SubscriptionID uuid.UUID) error {
	_, err := s.repo.Get(SubscriptionID)
	if err != nil {
		return fmt.Errorf("could not find lesson with ID %s: %v", SubscriptionID, err)
	}

	if err := s.repo.Delete(SubscriptionID); err != nil {
		return fmt.Errorf("failed to delete lesson with ID %s: %v", SubscriptionID, err)
	}

	log.Printf("Successfully deleted lesson with ID %s", SubscriptionID)
	return nil
}

func  (s *SubscriptionServiceImpl)GetAllSubscription()([]*model.Subscription, error) {
	 Subscription, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get all lesson: %v", err)
	}
	return Subscription, nil
}

func (s *SubscriptionServiceImpl) GetSubscriptionByID(SubscriptionID uuid.UUID) (*model.Subscription, error) {
	Subscription, err := s.repo.Get(SubscriptionID)
	if err != nil {
		return nil, err
	}
	return Subscription, nil
}

func (s *SubscriptionServiceImpl) UpdateSubscription(Subscription *model.Subscription) error {
  	_, err := s.repo.Get(Subscription.ID)
	if err != nil {
		return fmt.Errorf("could not find lesson with ID %s", Subscription.ID)
	}

	if err := s.repo.Update(Subscription); err != nil {
		return fmt.Errorf("failed to update lesson with ID %s: %v", Subscription.ID, err)
	}

	return nil
}


