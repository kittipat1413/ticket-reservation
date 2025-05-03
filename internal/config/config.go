package config

import (
	"fmt"
	"log"

	cfgFramework "github.com/kittipat1413/go-common/framework/config"
)

type Config struct {
	App     AppConfig            // Application-level settings such as API keys or feature flags
	Service ServiceConfig        // Infrastructure-level service settings like name, port, and environment
	DB      DatabaseConfig       // Database connection and pooling configuration
	Source  *cfgFramework.Config // Underlying unstructured config source, used for accessing unmapped keys
}

func MustConfigure() *Config {
	cfg, err := Configure()
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}

func Configure() (*Config, error) {
	cfg := cfgFramework.MustConfig(
		cfgFramework.WithOptionalConfigPaths("env.yaml", "../env.yaml"),
		cfgFramework.WithDefaults(configDefaults),
	)

	dbCfg, err := LoadDatabaseConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid database config: %w", err)
	}

	return &Config{
		App:     LoadAppConfig(cfg),
		Service: LoadServiceConfig(cfg),
		DB:      *dbCfg,
		Source:  cfg,
	}, nil
}
