package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Viper *viper.Viper
}

// Service configuration environment variable keys
const (
	serviceNameKey      = "SERVICE_NAME"
	servicePortKey      = "SERVICE_PORT"
	serviceEnvKey       = "SERVICE_ENV"
	serviceErrPrefixKey = "SERVICE_ERR_PREFIX"
	otelExporterKey     = "OTEL_EXPORTER_OTLP_ENDPOINT"
	adminApiKey         = "ADMIN_API_KEY"
	adminApiSecret      = "ADMIN_API_SECRET"
)

// Database configuration keys
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// Database configuration environment variable keys
const (
	databaseUrlKey             = "DATABASE_URL"
	databaseMaxOpenConnsKey    = "DATABASE_MAX_OPEN_CONNS"
	databaseMaxIdleConnsKey    = "DATABASE_MAX_IDLE_CONNS"
	databaseConnMaxLifetimeKey = "DATABASE_CONN_MAX_LIFETIME"  // duration string like "30m"
	databaseConnMaxIdleTimeKey = "DATABASE_CONN_MAX_IDLE_TIME" // duration string like "5m"
)

func MustConfigure() *Config {
	if cfg, err := Configure(); err != nil {
		log.Fatalln(err)
		return nil
	} else {
		return cfg
	}
}

func Configure() (*Config, error) {
	v, err := initViper()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}
	return &Config{
		Viper: v,
	}, nil
}

func initViper() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()
	if _, err := os.Stat("env.yaml"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile("env.yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	} else if _, err := os.Stat("../env.yaml"); !errors.Is(err, os.ErrNotExist) {
		v.SetConfigFile("../env.yaml")
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	v.SetDefault(serviceNameKey, "ticket-reservation-api")
	v.SetDefault(servicePortKey, ":8080")
	v.SetDefault(serviceEnvKey, "development")
	v.SetDefault(serviceErrPrefixKey, "TR")
	v.SetDefault(databaseUrlKey, "postgres:///ticket-reservation?sslmode=disable")
	v.SetDefault(databaseMaxOpenConnsKey, 30)
	v.SetDefault(databaseMaxIdleConnsKey, 15)
	v.SetDefault(databaseConnMaxLifetimeKey, "30m")
	v.SetDefault(databaseConnMaxIdleTimeKey, "5m")
	return v, nil
}

func (c *Config) ServiceName() string      { return c.Viper.GetString(serviceNameKey) }
func (c *Config) ServicePort() string      { return c.Viper.GetString(servicePortKey) }
func (c *Config) ServiceEnv() string       { return c.Viper.GetString(serviceEnvKey) }
func (c *Config) ServiceErrPrefix() string { return c.Viper.GetString(serviceErrPrefixKey) }
func (c *Config) OtelExporter() string     { return c.Viper.GetString(otelExporterKey) }
func (c *Config) AdminApiKey() string      { return c.Viper.GetString(adminApiKey) }
func (c *Config) AdminApiSecret() string   { return c.Viper.GetString(adminApiSecret) }
func (c *Config) DatabaseURL() string      { return c.Viper.GetString(databaseUrlKey) }
func (c *Config) DatabaseConfig() (*DatabaseConfig, error) {
	lifetime, err := time.ParseDuration(c.Viper.GetString(databaseConnMaxLifetimeKey))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_CONN_MAX_LIFETIME: %w", err)
	}

	idleTime, err := time.ParseDuration(c.Viper.GetString(databaseConnMaxIdleTimeKey))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_CONN_MAX_IDLE_TIME: %w", err)
	}

	return &DatabaseConfig{
		URL:             c.Viper.GetString(databaseUrlKey),
		MaxOpenConns:    c.Viper.GetInt(databaseMaxOpenConnsKey),
		MaxIdleConns:    c.Viper.GetInt(databaseMaxIdleConnsKey),
		ConnMaxLifetime: lifetime,
		ConnMaxIdleTime: idleTime,
	}, nil
}

func (c *Config) AllConfigurations() map[string]interface{} {
	m := map[string]interface{}{}
	for _, key := range c.Viper.AllKeys() {
		m[key] = c.Viper.Get(key)
	}
	return m
}
