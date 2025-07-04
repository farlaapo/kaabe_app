package gateway

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"

	"github.com/gofrs/uuid"
)

type PaymentRepositoryImpl struct {
	db *sql.DB
}

// Create implements repository.PaymentRepository.
func (p *PaymentRepositoryImpl) Create(payment *model.Payment) error {
	_, err := p.db.Exec(`call create_payment($1, $2, $3, $4, $5, $6, $7)`,
		payment.ID, payment.ExternalRef, payment.UserID, payment.SubscriptionID, payment.Amount, payment.Status, payment.ProcessedAt)

	//
	if err != nil {
		log.Printf("call create_payment error: %v", err)
		return err
	}

	log.Printf("payment created : %v", payment)
	return nil
}

// Delete implements repository.PaymentRepository.
func (p *PaymentRepositoryImpl) Delete(paymentID uuid.UUID) error {
	//
	_, err := p.db.Exec(`call delete_payment($1)`, paymentID)
	if err != nil {
		log.Printf("call delete _payment error: %v", err)
		return err
	}
	log.Printf("delete payment : %v", paymentID)
	return nil
}

// GetAll implements repository.PaymentRepository.
func (p *PaymentRepositoryImpl) GetAll() ([]*model.Payment, error) {
	//
	rows, err := p.db.Query(`select * from get_all_payments()`)
	if err != nil {
		log.Printf("call get_all_payment error: %v", err)
		return nil, err
	}

	defer rows.Close()
	var payments []*model.Payment

	for rows.Next() {
		var payment model.Payment
		err = rows.Scan(
			&payment.ID,
			&payment.ExternalRef,
			&payment.UserID,
			&payment.SubscriptionID,
			&payment.Amount,
			&payment.Status,
			&payment.ProcessedAt,
			&payment.UpdatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			log.Printf("scan error payment row: %v", err)
			return nil, err
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		log.Printf(" Row iteration error: %v", err)
		return nil, err
	}
	log.Printf("get all payment : %v", payments)
	return payments, nil
}

// GetByID implements repository.PaymentRepository.
func (p *PaymentRepositoryImpl) GetByID(paymentID uuid.UUID) (*model.Payment, error) {
	//
	var payment model.Payment
	row := p.db.QueryRow(`select * from get_payment_by_id($1)`, paymentID)

	err := row.Scan(
		&payment.ID,
		&payment.ExternalRef,
		&payment.UserID,
		&payment.SubscriptionID,
		&payment.Amount,
		&payment.Status,
		&payment.ProcessedAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("paymnet not found with ID: %v", paymentID)
			return nil, fmt.Errorf("payment not found")
		}
		log.Printf("Error scanning rating by ID: %v", err)
		return nil, err
	}

	log.Printf("payment retrieved by ID: %+v", payment)
	return &payment, nil

}

// Update implements repository.PaymentRepository.
func (p *PaymentRepositoryImpl) Update(payment *model.Payment) error {
	//
	_, err := p.db.Exec(`CALL update_payment($1, $2, $3, $4, $5, $6, $7)`,
		payment.ID, payment.ExternalRef, payment.UserID, payment.SubscriptionID, payment.Amount, payment.Status, payment.ProcessedAt)
	if err != nil {
		log.Printf("Error calling update_payment: %v", err)
		return err
	}

	log.Printf("payment updated: %+v", payment)
	return nil

}


func (p *PaymentRepositoryImpl) GetByExternalRef(ref string) (*model.Payment, error) {
	var payment model.Payment
	row := p.db.QueryRow(`SELECT * FROM payments WHERE external_ref = $1 AND deleted_at IS NULL`, ref)
	err := row.Scan(&payment.ID, &payment.ExternalRef, &payment.UserID, &payment.SubscriptionID, &payment.Amount, &payment.Status, &payment.ProcessedAt, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// NewPaymentRepositoryImpl returns a new instance of PaymentRepositoryImpl
func NewPaymentRepository(db *sql.DB) repository.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}
