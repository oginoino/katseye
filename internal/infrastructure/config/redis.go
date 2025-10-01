package config

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
	rediscache "katseye/internal/infrastructure/cache/redis"
)

type RedisResources struct {
	Client *goredis.Client
	TTL    time.Duration
}

func newRedisResources(cfg CacheConfig) (*RedisResources, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	client, err := rediscache.NewClient(rediscache.Config{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return nil, err
	}

	ttl := cfg.Redis.TTL
	if ttl <= 0 {
		ttl = defaultRedisTTL
	}

	return &RedisResources{
		Client: client,
		TTL:    ttl,
	}, nil
}

func (r *RedisResources) Close(ctx context.Context) error {
	if r == nil || r.Client == nil {
		return nil
	}

	return r.Client.Close()
}
