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

func (r *redisOTPRepo) StoreCode(target string, code string, expiresAt time.Time) error {
	ctx := context.Background()
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = defaultCodeTTL
	}
	return r.redis.Client.Set(ctx, "otp:"+target, code, ttl).Err()
}

func (r *redisOTPRepo) GetCode(target string) (code string, expiresAt time.Time, found bool, err error) {
	ctx := context.Background()
	val, err := r.redis.Client.Get(ctx, "otp:"+target).Result()
	if err != nil {
		return "", time.Time{}, false, nil // Not found
	}
	ttl, _ := r.redis.Client.TTL(ctx, "otp:"+target).Result()
	return val, time.Now().Add(ttl), true, nil
}

func (r *redisOTPRepo) DeleteCode(target string) error {
	return r.redis.Client.Del(context.Background(), "otp:"+target).Err()
}
