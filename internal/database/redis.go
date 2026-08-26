package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"SDOBA/internal/config"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf(
		"%s:%s",
		cfg.Redis.Host,
		cfg.Redis.Port,
	)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
