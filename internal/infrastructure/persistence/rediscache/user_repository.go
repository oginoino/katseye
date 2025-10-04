package rediscache

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"katseye/internal/domain/entities"
	"katseye/internal/domain/repositories"
)

type userRepository struct {
	repo   repositories.UserRepository
	client *goredis.Client
	ttl    time.Duration
}

type cachedUser struct {
	ID           primitive.ObjectID       `json:"id"`
	Email        string                   `json:"email"`
	PasswordHash string                   `json:"password_hash"`
	Active       bool                     `json:"active"`
	Role         entities.Role            `json:"role"`
	Permissions  []string                 `json:"permissions,omitempty"`
	ProfileType  entities.UserProfileType `json:"profile_type"`
	ProfileID    primitive.ObjectID       `json:"profile_id"`
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
		var payload cachedUser
		if decodeErr := json.Unmarshal(data, &payload); decodeErr == nil {
			user := payload.toEntity()
			if user != nil && user.PasswordHash != "" {
				log.Printf("cache: debug password_hash_len=%d prefix=%q", len(user.PasswordHash), prefix(user.PasswordHash, 7))
				log.Printf("cache: hit resource=users operation=find_by_email email=%s source=redis", normalized)
				return user, nil
			}
			log.Printf("cache: stale resource=users operation=find_by_email email=%s reason=missing_password_hash", normalized)
			if !payload.ID.IsZero() {
				r.evictUser(ctx, payload.ID)
			} else {
				_ = r.client.Del(ctx, key).Err()
			}
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
		var payload cachedUser
		if decodeErr := json.Unmarshal(data, &payload); decodeErr == nil {
			user := payload.toEntity()
			if user != nil && user.PasswordHash != "" {
				log.Printf("cache: hit resource=users operation=find_by_id id=%s source=redis", id.Hex())
				return user, nil
			}
			log.Printf("cache: stale resource=users operation=find_by_id id=%s reason=missing_password_hash", id.Hex())
			_ = r.client.Del(ctx, key).Err()
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

	cached := newCachedUser(user)
	if cached == nil {
		return
	}

	payload, err := json.Marshal(cached)
	if err != nil {
		return
	}

	log.Printf("cache: debug writing password_hash_len=%d prefix=%q", len(cached.PasswordHash), prefix(cached.PasswordHash, 7))

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
	var payload cachedUser
	if data, err := r.client.Get(ctx, key).Bytes(); err == nil {
		if unmarshalErr := json.Unmarshal(data, &payload); unmarshalErr == nil {
			email := strings.TrimSpace(strings.ToLower(payload.Email))
			if email != "" {
				_ = r.client.Del(ctx, buildEmailKey(email)).Err()
			}
		}
	}
	_ = r.client.Del(ctx, key).Err()
}

func newCachedUser(user *entities.User) *cachedUser {
	if user == nil {
		return nil
	}

	perms := make([]string, len(user.Permissions))
	copy(perms, user.Permissions)

	return &cachedUser{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Active:       user.Active,
		Role:         user.Role,
		Permissions:  perms,
		ProfileType:  user.ProfileType,
		ProfileID:    user.ProfileID,
	}
}

func (u *cachedUser) toEntity() *entities.User {
	if u == nil {
		return nil
	}

	perms := make([]string, len(u.Permissions))
	copy(perms, u.Permissions)

	user := &entities.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Active:       u.Active,
		Role:         u.Role,
		Permissions:  perms,
		ProfileType:  u.ProfileType,
		ProfileID:    u.ProfileID,
	}

	user.Normalize()
	return user
}

func prefix(v string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(v) <= n {
		return v
	}
	return v[:n]
}
