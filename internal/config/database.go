package config

import (
	"fmt"
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

func LoadDatabaseConfig(cfg *cfgFramework.Config) (*DatabaseConfig, error) {
	lifetime, err := time.ParseDuration(cfg.GetString(DatabaseConnMaxLifetimeKey))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_CONN_MAX_LIFETIME: %w", err)
	}

	idleTime, err := time.ParseDuration(cfg.GetString(DatabaseConnMaxIdleTimeKey))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_CONN_MAX_IDLE_TIME: %w", err)
	}

	return &DatabaseConfig{
		URL:             cfg.GetString(DatabaseUrlKey),
		MaxOpenConns:    cfg.GetInt(DatabaseMaxOpenConnsKey),
		MaxIdleConns:    cfg.GetInt(DatabaseMaxIdleConnsKey),
		ConnMaxLifetime: lifetime,
		ConnMaxIdleTime: idleTime,
	}, nil
}
