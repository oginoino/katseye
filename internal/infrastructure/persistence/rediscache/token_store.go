package rediscache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"katseye/internal/domain/security"
)

const (
	tokenRevocationNamespace = "auth:revoked:"
	minimumRevocationTTL     = time.Minute
)

var _ security.TokenStore = (*TokenStore)(nil)

// TokenStore persists revoked token metadata into Redis.
type TokenStore struct {
	client *goredis.Client
}

// NewTokenStore creates a TokenStore backed by the provided Redis client.
func NewTokenStore(client *goredis.Client) *TokenStore {
	if client == nil {
		return nil
	}

	return &TokenStore{client: client}
}

// Revoke stores the token hash with a TTL matching the remaining token lifetime.
func (s *TokenStore) Revoke(ctx context.Context, token string, expiresAt time.Time) error {
	if s == nil || s.client == nil {
		return nil
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}

	// Ensure we never store entries with a non-positive TTL to avoid immediate eviction.
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = minimumRevocationTTL
	}

	return s.client.Set(ctx, revocationKey(token), "revoked", ttl).Err()
}

// IsRevoked checks whether the token hash exists in Redis.
func (s *TokenStore) IsRevoked(ctx context.Context, token string) (bool, error) {
	if s == nil || s.client == nil {
		return false, nil
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return false, nil
	}

	exists, err := s.client.Exists(ctx, revocationKey(token)).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func revocationKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return tokenRevocationNamespace + hex.EncodeToString(hash[:])
}
