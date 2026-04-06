package otp

import (
	"context"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/data"
)

type redisOTPRepo struct {
	redis *data.RedisClient
}

func NewRedisOTPRepo(redis *data.RedisClient) otpRepo {
	return &redisOTPRepo{redis: redis}
}

func (r *redisOTPRepo) StoreCode(purpose string, target string, code string, expiresAt time.Time) error {
	ctx := context.Background()
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = defaultCodeTTL
	}
	key := "otp:" + purpose + ":" + target
	return r.redis.Client.Set(ctx, key, code, ttl).Err()
}

func (r *redisOTPRepo) GetCode(purpose string, target string) (code string, expiresAt time.Time, found bool, err error) {
	ctx := context.Background()
	key := "otp:" + purpose + ":" + target
	val, err := r.redis.Client.Get(ctx, key).Result()
	if err != nil {
		return "", time.Time{}, false, nil // Not found
	}
	ttl, _ := r.redis.Client.TTL(ctx, key).Result()
	return val, time.Now().Add(ttl), true, nil
}

func (r *redisOTPRepo) DeleteCode(purpose string, target string) error {
	key := "otp:" + purpose + ":" + target
	return r.redis.Client.Del(context.Background(), key).Err()
}
