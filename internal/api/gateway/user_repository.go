package gateway

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type userRepositoryImpl struct {
	db *sql.DB
}

// Create a new user
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
		user.WalletID,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("Error calling create_user: %v", err)
		return err
	}

	log.Printf("User created with ID: %s", user.ID)
	return nil
}

// Update user info
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
		user.UpdatedAt,
	).Scan(&updatedAt)

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	user.UpdatedAt = updatedAt
	log.Printf("User updated at: %v", updatedAt)
	return nil
}

// Delete user by ID
func (r *userRepositoryImpl) Delete(userID uuid.UUID) error {
	var rowsDeleted int

	err := r.db.QueryRow("SELECT delete_user($1)", userID).Scan(&rowsDeleted)
	if err != nil {
		log.Printf("Error calling delete_user: %v", err)
		return err
	}

	if rowsDeleted == 0 {
		log.Printf("User not found")
		return fmt.Errorf("user not found")
	}

	log.Printf("User deleted")
	return nil
}

// Get user by ID
func (r *userRepositoryImpl) Get(userID uuid.UUID) (*model.User, error) {
	var user model.User

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
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found")
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("DB error: %v", err)
		return nil, err
	}

	return &user, nil
}

// Find user by email
func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	var walletID uuid.NullUUID

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

	if walletID.Valid {
		wid := walletID.UUID.String()
		user.WalletID = &wid
	} else {
		user.WalletID = nil
	}

	return &user, nil
}


// List all users
func (r *userRepositoryImpl) List() ([]*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM get_all_users()")
	if err != nil {
		log.Printf("Error getting users: %v", err)
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
		log.Printf("Error with user rows: %v", err)
		return nil, err
	}

	return users, nil
}

// SetResetToken saves reset token and expiry
func (r *userRepositoryImpl) SetResetToken(email string, token uuid.UUID, expiry string) error {
	query := `SELECT set_reset_token($1, $2, $3)`
	_, err := r.db.Exec(query, email, token, expiry)
	if err != nil {
		log.Printf("Error setting reset token: %v", err)
		return fmt.Errorf("failed to set reset token: %w", err)
	}
	return nil
}

// FindByResetToken finds user by reset token
func (r *userRepositoryImpl) FindByResetToken(token uuid.UUID) (*model.User, error) {
	var user model.User
	var walletID uuid.NullUUID

	query := `SELECT * FROM get_user_by_reset_token($1)`
	err := r.db.QueryRow(query, token).Scan(
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
		if err == sql.ErrNoRows {
			log.Printf("Reset token invalid or expired")
			return nil, fmt.Errorf("reset token invalid or expired")
		}
		log.Printf("DB error: %v", err)
		return nil, fmt.Errorf("failed to find user by reset token: %w", err)
	}

	if walletID.Valid {
		wid := walletID.UUID.String()
		user.WalletID = &wid
	} else {
		user.WalletID = nil
	}

	return &user, nil
}

// UpdatePassword updates the user's hashed password
func (r *userRepositoryImpl) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
	query := `SELECT update_user_password($1, $2)`
	_, err := r.db.Exec(query, userID, hashedPassword)
	if err != nil {
		log.Printf("Error updating password: %v", err)
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// ClearResetToken clears the reset token fields
func (r *userRepositoryImpl) ClearResetToken(userID uuid.UUID) error {
	query := `SELECT clear_reset_token($1)`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		log.Printf("Error clearing reset token: %v", err)
		return fmt.Errorf("failed to clear reset token: %w", err)
	}
	return nil
}

// Factory
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}
