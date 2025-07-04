package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)



// Repository
type WithdrawalRepository interface {
	Create(Withdrawal *model.Withdrawal) error 
	Update(Withdrawal  *model.Withdrawal) error
	Delete(WithdrawalID uuid.UUID) error
	Get(WithdrawalID uuid.UUID) (*model.Withdrawal, error)
	List() ([]*model.Withdrawal, error)
}