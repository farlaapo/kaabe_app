package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

// Repository 
type SubscriptionRepository interface {
	Create(Subscription *model.Subscription) error 
	Update(Subscription  *model.Subscription) error
	Delete(SubscriptionID uuid.UUID) error
	Get(SubscriptionID uuid.UUID) (*model.Subscription, error)
	List() ([]*model.Subscription, error)

}