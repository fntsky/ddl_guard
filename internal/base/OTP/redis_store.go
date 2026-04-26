package otp

import (
	"context"
	"errors"
	"time"

	"github.com/fntsky/ddl_guard/internal/base/data"
	apperrors "github.com/fntsky/ddl_guard/internal/errors"
	"github.com/redis/go-redis/v9"
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
	if err := r.redis.Client.Set(ctx, key, code, ttl).Err(); err != nil {
		return apperrors.RedisError("otp.store", err)
	}
	return nil
}

func (r *redisOTPRepo) GetCode(purpose string, target string) (code string, expiresAt time.Time, found bool, err error) {
	ctx := context.Background()
	key := "otp:" + purpose + ":" + target

	val, err := r.redis.Client.Get(ctx, key).Result()
	if err != nil {
		// key 不存在是预期行为，不是错误
		if errors.Is(err, redis.Nil) {
			return "", time.Time{}, false, nil
		}
		return "", time.Time{}, false, apperrors.RedisError("otp.get", err)
	}

	ttl, err := r.redis.Client.TTL(ctx, key).Result()
	if err != nil {
		return "", time.Time{}, false, apperrors.RedisError("otp.ttl", err)
	}

	return val, time.Now().Add(ttl), true, nil
}

func (r *redisOTPRepo) DeleteCode(purpose string, target string) error {
	key := "otp:" + purpose + ":" + target
	if err := r.redis.Client.Del(context.Background(), key).Err(); err != nil {
		return apperrors.RedisError("otp.delete", err)
	}
	return nil
}