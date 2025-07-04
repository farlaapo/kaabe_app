package gateway

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"fmt"
)

// SubscriptionImpl repositry
type SubscriptionImpl struct {
	db *sql.DB
}

// Create implements repository.SubscriptionRepository.
func (r *SubscriptionImpl) Create(Subscription *model.Subscription) error {
	_, err := r.db.Exec(`CALL create_subscription($1, $2, $3, $4, $5, $6)`,
		Subscription.ID, Subscription.UserID, Subscription.CourseID, Subscription.StartedAt, Subscription.ExpiresAt, Subscription.Status)
	if err != nil {
		log.Printf("Error calling create_subscriptio: %v", err)
		return err
	}

	log.Printf("Subscription created: %+v", Subscription)
	return nil
}

// Delete implements repository.SubscriptionRepository.
func (r *SubscriptionImpl) Delete(SubscriptionID uuid.UUID) error {
	_, err := r.db.Exec(`CALL  delete_subscription($1)`, SubscriptionID)
	if err != nil {
		log.Printf("Error calling delete_subscription for ID %v: %v", SubscriptionID, err)
		return err
	}

	log.Printf("Subscription soft-deleted: %v", SubscriptionID)
	return nil
}

// Get implements repository.SubscriptionRepository.
func (r *SubscriptionImpl) Get(SubscriptionID uuid.UUID) (*model.Subscription, error) {
	var Subscription model.Subscription

	row := r.db.QueryRow(`SELECT * FROM get_subscription_by_id($1)`, SubscriptionID)

	err := row.Scan(
		&Subscription.ID, 
		&Subscription.UserID,
		&Subscription.CourseID,
		&Subscription.StartedAt,
		&Subscription.ExpiresAt,
		&Subscription.Status,
		&Subscription.CreatedAt,
		&Subscription.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Subscription not found with ID: %v", SubscriptionID)
			return nil, fmt.Errorf("Subscription not found")
		}
		log.Printf("Error scanning Subscription by ID: %v", err)
		return nil, err
	}

	log.Printf("Subscription retrieved by ID: %+v", Subscription)
	return &Subscription, nil
}

// List implements repository.SubscriptionRepository.
func (r *SubscriptionImpl) List() ([]*model.Subscription, error) {
	rows, err := r.db.Query(`SELECT * FROM get_all_subscriptions()`)
	if err != nil {
		log.Printf("Error querying get_all_Subscription: %v", err)
		return nil, err
	}
	defer rows.Close()

	var Subscriptions []*model.Subscription

	for rows.Next() {
		var Subscription model.Subscription
		err := rows.Scan(
			&Subscription.ID,
			&Subscription.UserID,
			&Subscription.CourseID,
			&Subscription.StartedAt,
			&Subscription.ExpiresAt,
			&Subscription.Status,
			&Subscription.CreatedAt,
			&Subscription.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning Subscription row: %v", err)
			return nil, err
		}
		Subscriptions = append(Subscriptions, &Subscription)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("Subscription retrieved: %d", len(Subscriptions))
	return Subscriptions, nil
}

// Update implements repository.SubscriptionRepository.
func (r *SubscriptionImpl) Update(Subscription *model.Subscription) error {
	_, err := r.db.Exec(`CALL update_subscription($1, $2, $3, $4)`,
		Subscription.ID,
		Subscription.StartedAt,
		Subscription.ExpiresAt,
		Subscription.Status, // ✅ FIXED
	)
	if err != nil {
		log.Printf("Error calling update_subscription: %v", err)
		return err
	}

	log.Printf("Subscription updated: %+v", Subscription)
	return nil
}


func NewSubscriptionImpl(db *sql.DB) repository.SubscriptionRepository {
	return &SubscriptionImpl{db: db}
}
