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
	// Create a new user
	RegisterUser(email, password, firstName, lastName, role, walletID string) (*model.User, error)

	// Authenticate user
	AuthenticateUser(email, password string) (*model.User, error)

	// Get user by ID
	GetUserByID(userID uuid.UUID) (*model.User, error)

	// Update user
	UpdateUser(user *model.User) error

	// Delete user
	DeleteUser(userID uuid.UUID) error

	// List all users
	ListUsers() ([]*model.User, error)

	// Forgot/reset password
	ForgotPassword(email string) error
	ResetPassword(token uuid.UUID, newPassword string) error
}

type userService struct {
	repo      repository.UserRepository
	tokenRepo repository.TokenRepository
}

// Register a new user
func (s *userService) RegisterUser(email, password, firstName, lastName, role, walletID string) (*model.User, error) {
	// Check if user already exists
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %v", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        newID,
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		WalletID:  &walletID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate a user
func (s *userService) AuthenticateUser(email, password string) (*model.User, error) {
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Check the password hash
    if !utils.CheckPasswordHash(password, user.Password) {
        return nil, errors.New("invalid email or password")
    }

    // Validate role from DB
    validRoles := map[string]bool{
        "user":       true,
        "admin":      true,
        "influencer": true,
    }
    if !validRoles[user.Role] {
        user.Role = "user"
    }

    // Create token (optional session feature)
    newToken, err := uuid.NewV4()
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    token := &model.Token{
        ID:        newToken,
        UserID:    user.ID,
        Token:     newToken.String(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    if err := s.tokenRepo.Create(token); err != nil {
        return nil, errors.New("failed to save token")
    }

    return user, nil
}


// Get user by ID
func (s *userService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	user, err := s.repo.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

// Update user
func (s *userService) UpdateUser(user *model.User) error {
	_, err := s.repo.Get(user.ID)
	if err != nil {
		return errors.New("user not found")
	}
	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// Delete user
func (s *userService) DeleteUser(userID uuid.UUID) error {
	_, err := s.repo.Get(userID)
	if err != nil {
		log.Printf("User not found: %v", err)
		return fmt.Errorf("user not found")
	}
	if err := s.repo.Delete(userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	log.Printf("User deleted successfully")
	return nil
}

// List users
func (s *userService) ListUsers() ([]*model.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	return users, nil
}

// ForgotPassword sets a reset token and expiry for the user
func (s *userService) ForgotPassword(email string) error {
	// Just check if the user exists
	if _, err := s.repo.FindByEmail(email); err != nil {
		return errors.New("email not found")
	}

	// Generate reset token
	resetToken, err := uuid.NewV4()
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	expiry := time.Now().Add(15 * time.Minute).Format(time.RFC3339)

	// Store reset token
	if err := s.repo.SetResetToken(email, resetToken, expiry); err != nil {
		return fmt.Errorf("failed to store reset token: %v", err)
	}

	// 🚀 Simulate email sending (for now)
	log.Printf("Reset token sent to %s: %s", email, resetToken.String())

	return nil
}


// ResetPassword updates the user's password using a valid reset token
func (s *userService) ResetPassword(token uuid.UUID, newPassword string) error {
	user, err := s.repo.FindByResetToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	if err := s.repo.UpdatePassword(user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	if err := s.repo.ClearResetToken(user.ID); err != nil {
		return fmt.Errorf("failed to clear reset token: %v", err)
	}

	log.Printf("Password reset successful for user: %s", user.Email)
	return nil
}

// Factory
func NewUserService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) UserService {
	return &userService{
		repo:      userRepo,
		tokenRepo: tokenRepo,
	}
}
