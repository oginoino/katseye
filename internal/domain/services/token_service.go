package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"katseye/internal/domain/security"
)

// ErrTokenStoreUnavailable indicates the token store dependency was not configured.
var ErrTokenStoreUnavailable = errors.New("token store unavailable")

// TokenService encapsulates token revocation operations.
type TokenService struct {
	store security.TokenStore
}

// NewTokenService creates a new TokenService.
func NewTokenService(store security.TokenStore) *TokenService {
	if store == nil {
		return nil
	}
	return &TokenService{store: store}
}

// RevokeToken stores the provided token identifier until the expiration timestamp.
func (s *TokenService) RevokeToken(ctx context.Context, token string, expiresAt time.Time) error {
	if s == nil || s.store == nil {
		return ErrTokenStoreUnavailable
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}

	return s.store.Revoke(ctx, token, expiresAt)
}

// IsTokenRevoked returns true when the token has been revoked.
func (s *TokenService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	if s == nil || s.store == nil {
		return false, ErrTokenStoreUnavailable
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return false, nil
	}

	return s.store.IsRevoked(ctx, token)
}
