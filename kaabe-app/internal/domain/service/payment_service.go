package service

import (
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type PaymentService interface {
	CreatePayment(externalRef string, userID, subscriptionID uuid.UUID, amount float64, status string, processedAt time.Time) (*model.Payment, error)
	UpdatePayment(payment *model.Payment) error
	DeletePayment(paymentID uuid.UUID) error
	GetPaymentByID(paymentID uuid.UUID) (*model.Payment, error)
	GetAllPayments() ([]*model.Payment, error)
	GetPaymentByExternalRef(ref string) (*model.Payment, error) 
}

// paymentServiceImpl is the implementation
type PaymentServiceImpl struct {
	repo repository.PaymentRepository
}

// CreatePayment implements PaymentService.
func (p *PaymentServiceImpl) CreatePayment(externalRef string, userID uuid.UUID, subscriptionID uuid.UUID, amount float64, status string, processedAt time.Time) (*model.Payment, error) {
	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %v", err)
	}

	payment := &model.Payment{
		ID:             newID,
		ExternalRef:    externalRef,
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Amount:         amount,
		Status:         status,
		ProcessedAt:    processedAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	log.Printf("Creating payment: %+v", payment)

	if err := p.repo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	return payment, nil
}

// DeletePayment implements PaymentService.
func (p *PaymentServiceImpl) DeletePayment(paymentID uuid.UUID) error {
	_, err := p.repo.GetByID(paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	if err := p.repo.Delete(paymentID); err != nil {
		return fmt.Errorf("failed to delete payment: %v", err)
	}

	log.Printf("Deleted payment ID: %s", paymentID)
	return nil
}

// GetAllPayment implements PaymentService.
func (p *PaymentServiceImpl) GetAllPayments() ([]*model.Payment, error) {
	payments, err := p.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %v", err)
	}
	return payments, nil
}

// GetPaymentByID implements PaymentService.
func (p *PaymentServiceImpl) GetPaymentByID(paymentID uuid.UUID) (*model.Payment, error) {
	payment, err := p.repo.GetByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return payment, nil
}

// UpdatePayment implements PaymentService.
func (p *PaymentServiceImpl) UpdatePayment(payment *model.Payment) error {
	_, err := p.repo.GetByID(payment.ID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	if err := p.repo.Update(payment); err != nil {
		return fmt.Errorf("failed to update payment: %v", err)
	}

	log.Printf("Updated payment: %+v", payment)
	return nil
}

func (p *PaymentServiceImpl) GetPaymentByExternalRef(ref string) (*model.Payment, error) {
	return p.repo.GetByExternalRef(ref)
}

// NewPaymentService initializes the service
func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &PaymentServiceImpl{
		repo: paymentRepo,
	}
}
