package services

import (
	"context"
	"errors"
	"strings"

	interfaces "katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"
)

var (
	// ErrInvalidCredentials indicates the provided login information is incorrect.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInactiveAccount indicates the user exists but is not active.
	ErrInactiveAccount = errors.New("inactive account")
)

// AuthService handles credential verification against persisted users.
type AuthService struct {
	userRepo interfaces.UserRepository
}

func NewAuthService(userRepo interfaces.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Authenticate validates credentials, returning the user on success.
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*entities.User, error) {
	if s == nil || s.userRepo == nil {
		return nil, ErrInvalidCredentials
	}

	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive() {
		return nil, ErrInactiveAccount
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
