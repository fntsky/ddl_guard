package data

import (
	"context"
	"fmt"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() (*RedisClient, func(), error) {
	cfg := conf.Global()
	if cfg == nil || cfg.Redis.Addr == "" {
		// Redis not configured, return nil without error
		return nil, func() {}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, func() {}, fmt.Errorf("redis ping failed: %w", err)
	}

	cleanup := func() {
		client.Close()
	}

	return &RedisClient{Client: client}, cleanup, nil
}
