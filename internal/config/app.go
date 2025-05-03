package config

import cfgFramework "github.com/kittipat1413/go-common/framework/config"

type AppConfig struct {
	AdminAPIKey    string
	AdminAPISecret string
	// Add business feature flags here
}

func LoadAppConfig(cfg *cfgFramework.Config) AppConfig {
	return AppConfig{
		AdminAPIKey:    cfg.GetString(AdminApiKey),
		AdminAPISecret: cfg.GetString(AdminApiSecret),
	}
}
