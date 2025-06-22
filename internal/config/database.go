package config

import (
	"time"

	cfgFramework "github.com/kittipat1413/go-common/framework/config"
)

// Database configuration keys
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func LoadDatabaseConfig(cfg *cfgFramework.Config) DatabaseConfig {
	return DatabaseConfig{
		URL:             cfg.GetString(DatabaseURLKey),
		MaxOpenConns:    cfg.GetInt(DatabaseMaxOpenConnsKey),
		MaxIdleConns:    cfg.GetInt(DatabaseMaxIdleConnsKey),
		ConnMaxLifetime: cfg.GetDuration(DatabaseConnMaxLifetimeKey),
		ConnMaxIdleTime: cfg.GetDuration(DatabaseConnMaxIdleTimeKey),
	}
}

// Redis configuration keys
type RedisConfig struct {
	Addrs        []string
	Username     string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadRedisConfig(cfg *cfgFramework.Config) RedisConfig {
	return RedisConfig{
		Addrs:        cfg.GetStringSlice(RedisAddrsKey),
		Username:     cfg.GetString(RedisUsernameKey),
		Password:     cfg.GetString(RedisPasswordKey),
		DB:           cfg.GetInt(RedisDBKey),
		DialTimeout:  cfg.GetDuration(RedisDialTimeoutKey),
		ReadTimeout:  cfg.GetDuration(RedisReadTimeoutKey),
		WriteTimeout: cfg.GetDuration(RedisWriteTimeoutKey),
	}
}
