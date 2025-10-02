package services

import (
	"context"
	"errors"
	"strings"

	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrInvalidCredentials indicates the provided login information is incorrect.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInactiveAccount indicates the user exists but is not active.
	ErrInactiveAccount = errors.New("inactive account")
	// ErrInvalidUserData indicates the payload for user management is incomplete or invalid.
	ErrInvalidUserData = errors.New("invalid user data")
	// ErrInvalidRole indicates an unsupported role assignment was requested.
	ErrInvalidRole = errors.New("invalid role")
	// ErrUserAlreadyExists indicates an attempt to create a user with a duplicated email.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound indicates lookups failed to locate the target user.
	ErrUserNotFound = errors.New("user not found")
)

// AuthService handles credential verification against persisted users.
type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
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

// CreateUser provisions a new authenticated user with the provided credentials and authorisation metadata.
func (s *AuthService) CreateUser(ctx context.Context, email, password string, active bool, role entities.Role, permissions []string) (*entities.User, error) {
	if s == nil || s.userRepo == nil {
		return nil, ErrInvalidUserData
	}

	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return nil, ErrInvalidUserData
	}

	if role == "" {
		role = entities.RoleUser
	}
	if !entities.IsValidRole(role) {
		return nil, ErrInvalidRole
	}

	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &entities.User{
		ID:          primitive.NewObjectID(),
		Email:       email,
		Active:      active,
		Role:        role,
		Permissions: permissions,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	user.Normalize()

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser removes a user identified by the provided ID.
func (s *AuthService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	if s == nil || s.userRepo == nil {
		return ErrInvalidUserData
	}
	if id.IsZero() {
		return ErrInvalidUserData
	}

	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

// GetUserByID retrieves a user by identifier.
func (s *AuthService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	if s == nil || s.userRepo == nil {
		return nil, ErrInvalidUserData
	}
	if id.IsZero() {
		return nil, ErrInvalidUserData
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}
