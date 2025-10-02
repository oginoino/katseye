package rediscache

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func TestTokenStore_RevokeAndCheck(t *testing.T) {
	ctx := context.Background()
	server := miniredis.RunT(t)
	client := goredis.NewClient(&goredis.Options{Addr: server.Addr()})

	store := NewTokenStore(client)
	if store == nil {
		t.Fatal("expected token store instance")
	}

	const token = "sample.jwt.token"

	if revoked, err := store.IsRevoked(ctx, token); err != nil {
		t.Fatalf("IsRevoked returned error: %v", err)
	} else if revoked {
		t.Fatalf("expected token not to be revoked")
	}

	expiresAt := time.Now().Add(time.Hour)
	if err := store.Revoke(ctx, token, expiresAt); err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}

	if revoked, err := store.IsRevoked(ctx, token); err != nil {
		t.Fatalf("IsRevoked after revoke returned error: %v", err)
	} else if !revoked {
		t.Fatalf("expected token to be marked revoked")
	}
}

func TestTokenStore_RevokeWithPastExpirationUsesMinimumTTL(t *testing.T) {
	ctx := context.Background()
	server := miniredis.RunT(t)
	client := goredis.NewClient(&goredis.Options{Addr: server.Addr()})

	store := NewTokenStore(client)
	const token = "expired.jwt.token"

	if err := store.Revoke(ctx, token, time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}

	if ttl, err := client.TTL(ctx, revocationKey(token)).Result(); err != nil {
		t.Fatalf("TTL returned error: %v", err)
	} else if ttl < minimumRevocationTTL-time.Minute || ttl > minimumRevocationTTL+time.Minute {
		t.Fatalf("expected TTL around %s, got %s", minimumRevocationTTL, ttl)
	}
}
