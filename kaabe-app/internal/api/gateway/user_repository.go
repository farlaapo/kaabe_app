package gateway

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"time"

	"log"

	"github.com/gofrs/uuid"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func (r *userRepositoryImpl) Create(user *model.User) error {
	validRoles := map[string]bool{"admin": true, "user": true, "influencer": true}
	if !validRoles[user.Role] {
		return fmt.Errorf("invalid user role: %s", user.Role)
	}

	query := `SELECT create_user($1, $2, $3, $4, $5, $6);`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.WalletID, // *uuid.UUID or nil
	).Scan(&user.ID)

	if err != nil {
		log.Printf("Error calling create_user: %v", err)
		return err
	}

	log.Printf("User created with ID: %s", user.ID)
	return nil
}

// Delete implements repository.UserRepository.
func (r *userRepositoryImpl) Delete(userID uuid.UUID) error {
	var rowsDeleted int

	// Call the function with SELECT and scan the result
	err := r.db.QueryRow("SELECT delete_user($1)", userID).Scan(&rowsDeleted)
	if err != nil {
		log.Printf("Error calling delete_user function: %v", err)
		return err
	}

	if rowsDeleted == 0 {
		log.Printf("User not found")
		return fmt.Errorf("user not found")
	}

	log.Printf("Rows affected: %v", rowsDeleted)
	return nil
}

// FindByEmail implements repository.UserRepository.
func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var (
		user     model.User
		walletID uuid.NullUUID
	)

	// Set a default valid role before calling the procedure (Postgres enum cannot be empty string)
	if user.Role == "" {
		user.Role = "user"
	}

	query := `SELECT * FROM get_user_by_email($1)`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&walletID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("DB error: %v", err)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Convert walletID from NullUUID to *string
	if walletID.Valid {
		wid := walletID.UUID.String()
		user.WalletID = &wid
	} else {
		user.WalletID = nil
	}

	return &user, nil
}

// Get implements repository.UserRepository.
func (r *userRepositoryImpl) Get(userID uuid.UUID) (*model.User, error) {
	// define user
	var user model.User

	// get user from database
	query := `SELECT * FROM get_user_by_id($1)`
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.WalletID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	// Check for errors in retrieving the users
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found")
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("DB error: %v", err)
	}

	// return the user if found
	return &user, nil

}

// List implements repository.UserRepository.
func (r *userRepositoryImpl) List() ([]*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM get_all_users()")
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}

	defer rows.Close()
	var users []*model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.WalletID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error getting users: %v", err)
		return nil, err
	}

	return users, nil
}

// Update implements repository.UserRepository.
func (r *userRepositoryImpl) Update(user *model.User) error {
	query := `CALL update_user($1, $2, $3, $4, $5, $6, $7, $8)`
	var updatedAt time.Time

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.WalletID,
		user.UpdatedAt, // input value
	).Scan(&updatedAt) // scan output value

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	user.UpdatedAt = updatedAt // update model with new timestamp
	log.Printf("User updated at: %v", updatedAt)
	return nil
}

// NewUserRepository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}

}
