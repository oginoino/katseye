package security

import (
	"context"
	"time"
)

// TokenStore provides access to persisted token revocation metadata.
type TokenStore interface {
	// Revoke marks the provided token as revoked until the supplied expiration time.
	Revoke(ctx context.Context, token string, expiresAt time.Time) error
	// IsRevoked reports whether the token has been revoked.
	IsRevoked(ctx context.Context, token string) (bool, error)
}

