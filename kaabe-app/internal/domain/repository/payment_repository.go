package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	Update(payment *model.Payment) error
	Delete(paymentID uuid.UUID) error
	GetByID(paymentID uuid.UUID) (*model.Payment, error)
	GetAll() ([]*model.Payment, error)
	GetByExternalRef(externalRef string) (*model.Payment, error)
}
