package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	Address  string
	Password string
	DB       int
}

func NewClient(cfg Config) (*goredis.Client, error) {
	options := &goredis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := goredis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
