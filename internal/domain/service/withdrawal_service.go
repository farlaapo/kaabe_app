package service

import (
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"

	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

// Service
type WithdrawalService interface {
	CreateWithdrawal(InfluencerID uuid.UUID, amount float64, status string, RequestedAt time.Time, ProcessedAt time.Time) (*model.Withdrawal, error)
	UpdateWithdrawal(withdrawal *model.Withdrawal) error
	DeleteWithdrawal(withdrawalID uuid.UUID) error
	GetWithdrawalByID(withdrawalID uuid.UUID) (*model.Withdrawal, error)
	GetAllWithdrawal() ([]*model.Withdrawal, error)
}

type withdrawalService struct {
	repo      repository.WithdrawalRepository
	tokenRepo repository.TokenRepository
}

// CreateWithdrawal implements WithdrawalService.
func (s *withdrawalService) CreateWithdrawal(InfluencerID uuid.UUID, amount float64, status string, RequestedAt time.Time, ProcessedAt time.Time) (*model.Withdrawal, error) {
		// Generate a new UUID for the lesson ID

	neoLesson, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Create a new lesson instance
	amWithdrawal := &model.Withdrawal{
		ID:        neoLesson,
		InfluencerID: InfluencerID,
		Amount: amount,
		Status:     status,
		RequestedAt: RequestedAt,
		ProcessedAt: ProcessedAt,

	}

	// Log the new lesson creation attempt
	log.Printf("Creating amWithdrawal: %+v", amWithdrawal)

	// Save the new lesson to the repository
	err = s.repo.Create(amWithdrawal)
	if err != nil {
		return nil, fmt.Errorf("failed to create amWithdrawal: %v", err)
	}

	return amWithdrawal, nil
}

// DeleteWithdrawal implements WithdrawalService.
func (s *withdrawalService) DeleteWithdrawal(withdrawalID uuid.UUID) error {
		_, err := s.repo.Get(withdrawalID)
	if err != nil {
		return fmt.Errorf("could not find withdrawal with ID %s: %v", withdrawalID, err)
	}

	if err := s.repo.Delete(withdrawalID); err != nil {
		return fmt.Errorf("failed to delete withdrawal with ID %s: %v", withdrawalID, err)
	}

	log.Printf("Successfully deleted withdrawal with ID %s", withdrawalID)
	return nil
}

// GetAllWithdrawal implements WithdrawalService.
func (s *withdrawalService) GetAllWithdrawal() ([]*model.Withdrawal, error) {
	Withdrawal, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get all leWithdrawalsson: %v", err)
	}
	return Withdrawal, nil
}

// GetWithdrawalByID implements WithdrawalService.
func (s *withdrawalService) GetWithdrawalByID(withdrawalID uuid.UUID) (*model.Withdrawal, error) {
	withdrawal, err := s.repo.Get(withdrawalID)
	if err != nil {
		return nil, err
	}
	return withdrawal, nil
}

// UpdateWithdrawal implements WithdrawalService.
func (s *withdrawalService) UpdateWithdrawal(withdrawal *model.Withdrawal) error {
	_, err := s.repo.Get(withdrawal.ID)
	if err != nil {
		return fmt.Errorf("could not find withdrawal with ID %s", withdrawal.ID)
	}

	if err := s.repo.Update(withdrawal); err != nil {
		return fmt.Errorf("failed to update withdrawal with ID %s: %v", withdrawal.ID, err)
	}

	return nil
}

func NewWithdrawalService(WithdrawalRepo repository.WithdrawalRepository, tokenRepo repository.TokenRepository) WithdrawalService {
	return &withdrawalService{
		repo:      WithdrawalRepo,
		tokenRepo: tokenRepo,
	}
}
