package rediscache

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"katseye/internal/application/interfaces/repositories"
	"katseye/internal/domain/entities"
)

type userRepository struct {
	repo   repositories.UserRepository
	client *goredis.Client
	ttl    time.Duration
}

func NewUserRepository(client *goredis.Client, ttl time.Duration, repo repositories.UserRepository) repositories.UserRepository {
	if client == nil || repo == nil {
		return repo
	}

	return &userRepository{
		repo:   repo,
		client: client,
		ttl:    mergeTTL(ttl, time.Minute),
	}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	if r == nil {
		return nil, nil
	}

	normalized := strings.TrimSpace(strings.ToLower(email))
	if normalized == "" {
		return r.repo.FindByEmail(ctx, email)
	}

	key := buildEmailKey(normalized)

	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached entities.User
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			log.Printf("cache: hit resource=users operation=find_by_email email=%s source=redis", normalized)
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=users operation=find_by_email email=%s error=%v", normalized, unmarshalErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		if payload, marshalErr := json.Marshal(user); marshalErr == nil {
			_ = r.client.Set(ctx, key, payload, r.ttl).Err()
		}
	}

	if user != nil {
		log.Printf("cache: miss resource=users operation=find_by_email email=%s source=mongo", normalized)
	} else {
		log.Printf("cache: miss resource=users operation=find_by_email email=%s source=mongo result=empty", normalized)
	}

	return user, nil
}

func buildEmailKey(email string) string {
	return "users:email:" + email
}
