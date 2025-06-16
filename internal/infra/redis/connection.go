package redis

import (
	"ticket-reservation/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.Config) redis.UniversalClient {
	options := &redis.UniversalOptions{
		Addrs:    cfg.Redis.Addrs,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	if cfg.Redis.ReadTimeout != 0 {
		options.ReadTimeout = cfg.Redis.ReadTimeout
	}
	if cfg.Redis.WriteTimeout != 0 {
		options.WriteTimeout = cfg.Redis.WriteTimeout
	}
	if cfg.Redis.DialTimeout != 0 {
		options.DialTimeout = cfg.Redis.DialTimeout
	}

	client := redis.NewUniversalClient(options)
	return client
}
