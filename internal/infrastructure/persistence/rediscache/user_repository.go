package rediscache

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		if decodeErr := json.Unmarshal(data, &cached); decodeErr == nil {
			log.Printf("cache: hit resource=users operation=find_by_email email=%s source=redis", normalized)
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=users operation=find_by_email email=%s error=%v", normalized, decodeErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		r.cacheUser(ctx, user)
	}

	if user != nil {
		log.Printf("cache: miss resource=users operation=find_by_email email=%s source=mongo", normalized)
	} else {
		log.Printf("cache: miss resource=users operation=find_by_email email=%s source=mongo result=empty", normalized)
	}

	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	if r == nil {
		return nil, nil
	}

	if id.IsZero() {
		return r.repo.FindByID(ctx, id)
	}

	key := buildIDKey("users", id.Hex())
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		var cached entities.User
		if decodeErr := json.Unmarshal(data, &cached); decodeErr == nil {
			log.Printf("cache: hit resource=users operation=find_by_id id=%s source=redis", id.Hex())
			return &cached, nil
		} else {
			log.Printf("cache: stale resource=users operation=find_by_id id=%s error=%v", id.Hex(), decodeErr)
			_ = r.client.Del(ctx, key).Err()
		}
	}

	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		r.cacheUser(ctx, user)
		log.Printf("cache: miss resource=users operation=find_by_id id=%s source=mongo", id.Hex())
	} else {
		log.Printf("cache: miss resource=users operation=find_by_id id=%s source=mongo result=empty", id.Hex())
	}

	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	if err := r.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	r.cacheUser(ctx, user)

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	if err := r.repo.DeleteUser(ctx, id); err != nil {
		return err
	}

	r.evictUser(ctx, id)

	return nil
}

func buildEmailKey(email string) string {
	return "users:email:" + email
}

func (r *userRepository) cacheUser(ctx context.Context, user *entities.User) {
	if r == nil || r.client == nil || user == nil {
		return
	}

	payload, err := json.Marshal(user)
	if err != nil {
		return
	}

	if !user.ID.IsZero() {
		_ = r.client.Set(ctx, buildIDKey("users", user.ID.Hex()), payload, r.ttl).Err()
	}

	email := strings.TrimSpace(strings.ToLower(user.Email))
	if email != "" {
		_ = r.client.Set(ctx, buildEmailKey(email), payload, r.ttl).Err()
	}
}

func (r *userRepository) evictUser(ctx context.Context, id primitive.ObjectID) {
	if r == nil || r.client == nil {
		return
	}

	key := buildIDKey("users", id.Hex())
	var cached entities.User
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		if unmarshalErr := json.Unmarshal(data, &cached); unmarshalErr == nil {
			email := strings.TrimSpace(strings.ToLower(cached.Email))
			if email != "" {
				_ = r.client.Del(ctx, buildEmailKey(email)).Err()
			}
		}
	}
	_ = r.client.Del(ctx, key).Err()
}
