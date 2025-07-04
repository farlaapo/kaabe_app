package service

import (
	"errors"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	utils "kaabe-app/pkg/config"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type UserService interface {

	// Create a New user
	RegisterUser(email, Password, FirstName, LastName, Role, WalletID string) (*model.User, error)
	// authentication
	AuthenticateUser(email, Password_hash string) (*model.User, error)
	// Get a User by id
	GetUserByID(userID uuid.UUID) (*model.User, error)
	// Update a user
	UpdateUser(user *model.User) error
	//Delete a user
	DeleteUser(userID uuid.UUID) error
	ListUsers() ([]*model.User, error)
}

type userService struct {
	repo      repository.UserRepository
	tokenRepo repository.TokenRepository
}

// RegisterUser implements UserService.
func (s *userService) RegisterUser(email string, Password string, FirstName string, LastName string, Role string, WalletID string) (*model.User, error) {
	// CHECK iF USER ALREADY EXISTS
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %v", err)
	}


	// Hash Password using bcrypt
	hashPassword, err := utils.HashPassword(Password)
	if err != nil {
		return nil, err
	}

	// Create a new user
	user := &model.User{
		ID: newID,
		Email:     email,
		Password:  hashPassword,
		FirstName: FirstName,
		LastName:  LastName,
		Role:      Role,
		WalletID:  &WalletID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// save user to database
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil

}


// AuthenticateUser implements UserService.
func (s *userService) AuthenticateUser(email string, password string) (*model.User, error) {
	// find email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// ✅ Validate role from DB (avoid invalid enum issue)
	validRoles := map[string]bool{
		"user":       true,
		"admin":      true,
		"influencer": true,
	}
	if !validRoles[user.Role] {
		user.Role = "user" // fallback to "user" if invalid
	}

	// generate token
	newToken, err := uuid.NewV4()
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Create token
	token := &model.Token{
		ID:        newToken,
		UserID:    user.ID,
		Token:     newToken.String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// store token in database
	if err := s.tokenRepo.Create(token); err != nil {
		return nil, errors.New("failed to save token")
	}

	return user, nil
}


// DeleteUser implements UserService.
func (s *userService) DeleteUser(userID uuid.UUID) error {
    // Check if user exists
    _, err := s.repo.Get(userID)
    if err != nil {
        log.Printf("User not found: %v", err)
        return fmt.Errorf("user not found")
    }

    // Delete user
    if err := s.repo.Delete(userID); err != nil {
        log.Printf("Error deleting user: %v", err)
        return err
    }

    log.Printf("User deleted successfully")
    return nil
}


// GetUserByID implements UserService.
func (s *userService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// ListUsers implements UserService.
func (s *userService) ListUsers() ([]*model.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	return users, nil
}

// UpdateUser implements UserService.
func (s *userService) UpdateUser(user *model.User) error {
	_, err := s.repo.Get(user.ID)
	if err != nil {
		return errors.New("user not found")
	}

	// call repository to update user
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) UserService {
	return &userService{
		repo:      userRepo,
		tokenRepo: tokenRepo,
	}
}
