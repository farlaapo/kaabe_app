package gateway

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"fmt"
)

type WithdrawalRepositoryImpl struct {
	db *sql.DB
}

// Create implements repository.WithdrawalRepository.
func (r *WithdrawalRepositoryImpl) Create(Withdrawal *model.Withdrawal) error {
	_, err := r.db.Exec(`CALL create_withdrawal($1, $2, $3, $4, $5, $6)`,
		Withdrawal.ID, Withdrawal.InfluencerID, Withdrawal.Amount, Withdrawal.Status, Withdrawal.RequestedAt, Withdrawal.ProcessedAt)
	if err != nil {
		log.Printf("Error calling create_withdrawal: %v", err)
		return err
	}

	log.Printf("Withdrawal created: %+v", Withdrawal)
	return nil
}

// Delete implements repository.WithdrawalRepository.
func (r *WithdrawalRepositoryImpl) Delete(WithdrawalID uuid.UUID) error {
	_, err := r.db.Exec(`CALL  delete_withdrawal($1)`, WithdrawalID)
	if err != nil {
		log.Printf("Error calling delete_withdrawal for ID %v: %v", WithdrawalID, err)
		return err
	}

	log.Printf("Withdrawal soft-deleted: %v", WithdrawalID)
	return nil
}

// Get implements repository.WithdrawalRepository.
func (r *WithdrawalRepositoryImpl) Get(WithdrawalID uuid.UUID) (*model.Withdrawal, error) {
	var Withdrawal model.Withdrawal

	row := r.db.QueryRow(`SELECT * FROM get_withdrawal_by_id($1)`, WithdrawalID)

	err := row.Scan(
		&Withdrawal.ID,
		&Withdrawal.InfluencerID,
		&Withdrawal.Amount,
		&Withdrawal.Status,
		&Withdrawal.RequestedAt,
		&Withdrawal.ProcessedAt,
		&Withdrawal.CreatedAt,
		&Withdrawal.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Withdrawal not found with ID: %v", WithdrawalID)
			return nil, fmt.Errorf("Withdrawal not found")
		}
		log.Printf("Error scanning Withdrawal by ID: %v", err)
		return nil, err
	}

	log.Printf("Withdrawal retrieved by ID: %+v", Withdrawal)
	return &Withdrawal, nil
}

// List implements repository.WithdrawalRepository.
func (r *WithdrawalRepositoryImpl) List() ([]*model.Withdrawal, error) {
	rows, err := r.db.Query(`SELECT * FROM get_all_withdrawals()`)
	if err != nil {
		log.Printf("Error querying get_all_withdrawals(): %v", err)
		return nil, err
	}
	defer rows.Close()

	var Withdrawals []*model.Withdrawal

	for rows.Next() {
		var Withdrawal model.Withdrawal
		err := rows.Scan(
			&Withdrawal.ID,
			&Withdrawal.InfluencerID,
			&Withdrawal.Amount,
			&Withdrawal.Status,
			&Withdrawal.RequestedAt,
			&Withdrawal.ProcessedAt,
			&Withdrawal.CreatedAt,
			&Withdrawal.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning Withdrawal row: %v", err)
			return nil, err
		}
		Withdrawals = append(Withdrawals, &Withdrawal)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("Withdrawal retrieved: %d", len(Withdrawals))
	return Withdrawals, nil
}

// Update implements repository.WithdrawalRepository.
func (r *WithdrawalRepositoryImpl) Update(Withdrawal *model.Withdrawal) error {
	_, err := r.db.Exec(`CALL update_withdrawal($1, $2, $3, $4)`,
		Withdrawal.ID,  Withdrawal.Amount, Withdrawal.Status, Withdrawal.ProcessedAt )
	if err != nil {
		log.Printf("Error calling update_withdrawal: %v", err)
		return err
	}

	log.Printf("Withdrawal updated: %+v", Withdrawal)
	return nil
}

func NewWithdrawalRepositoryImpl(db *sql.DB) repository.WithdrawalRepository {
	return &WithdrawalRepositoryImpl{db: db}

}
