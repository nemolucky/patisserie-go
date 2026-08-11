package redis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nemolucky/patisserie-go/internal/config"
	"github.com/nemolucky/patisserie-go/internal/domain"
)

type Redis struct {
	*redis.Client
	defaultExpiration time.Duration
}

func New(ctx context.Context, cfg config.RedisConfig, logger *slog.Logger) *Redis {

	rds := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})

	if _, err := rds.Ping(ctx).Result(); nil != err {
		logger.Error(
			"Redis connection error",
			slog.String("addr", cfg.Addr),
			slog.String("error", err.Error()),
		)
		return nil
	}
	return &Redis{
		rds,
		cfg.Expiration,
	}
}

func (rds *Redis) GetSMSCode(ctx context.Context, phone string) (string, error) {
	code, err := rds.Get(ctx, "sms:"+phone).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrNotFound
		}
		return "", err
	}

	return code, nil
}

func (rds *Redis) SetSMSCode(ctx context.Context, phone, code string) (error) {
	err := rds.Set(ctx, "sms:"+phone, code, rds.defaultExpiration).Err()
	if err != nil{
		return err
	}
	return nil
}

func (rds *Redis) DeleteSMSCode(ctx context.Context, phone string) (error) {
	err := rds.Del(ctx, "sms:"+phone).Err()
	if err != nil{
		return err
	}
	return nil
}

func (rds *Redis) SetVerified(ctx context.Context, phone string) (error) {
	err := rds.Set(ctx, "verified:"+phone, "1", rds.defaultExpiration).Err()
	if err != nil{
		return err
	}
	return nil
}

func (rds *Redis) IsVerified(ctx context.Context, phone string) (bool, error) {
	_, err := rds.Get(ctx, "verified:"+phone).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, domain.ErrNotFound
		}
		return false, err
	}

	return true, nil
}

func (rds *Redis) DeleteVerified(ctx context.Context, phone string) error {
	return rds.Del(ctx, "verified:"+phone).Err()
}