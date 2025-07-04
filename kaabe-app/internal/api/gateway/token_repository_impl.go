package gateway

import (
	"database/sql"
	"errors"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"time"
)

// tokenRepositoryImpl is the PostgreSQL-based implementation of TokenRepository
type tokenRepositoryImpl struct {
	db *sql.DB
}

// NewTokenRepository returns a new TokenRepository instance
func NewTokenRepository(db *sql.DB) repository.TokenRepository {
	return &tokenRepositoryImpl{db: db}
}

// Create inserts a new token into the database using a stored procedure
func (t *tokenRepositoryImpl) Create(token *model.Token) error {
	now := time.Now().UTC()
	token.CreatedAt = now
	token.UpdatedAt = now

	_, err := t.db.Exec(`
		CALL create_token($1, $2, $3, $4, $5, $6)
	`,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error executing stored procedure create_token: %v", err)
		return err
	}

	log.Printf("Token successfully created: ID=%v", token.ID)
	return nil
}

// FindByToken retrieves a token by its value using a SQL function
func (t *tokenRepositoryImpl) FindByToken(tokenStr string) (*model.Token, error) {
	query := `SELECT * FROM get_token_by_token($1);`

	row := t.db.QueryRow(query, tokenStr)

	var token model.Token
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.DeletedAt,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("Token not found: %s", tokenStr)
		return nil, nil
	case err != nil:
		log.Printf("Error scanning token row: %v", err)
		return nil, err
	default:
		return &token, nil
	}
}
