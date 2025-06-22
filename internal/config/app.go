package config

import (
	"time"

	cfgFramework "github.com/kittipat1413/go-common/framework/config"
)

type AppConfig struct {
	AdminAPIKey    string
	AdminAPISecret string
	Timezone       string
	SeatLockTTL    time.Duration
	// Add business feature flags here
}

func LoadAppConfig(cfg *cfgFramework.Config) AppConfig {
	return AppConfig{
		AdminAPIKey:    cfg.GetString(AdminAPIKey),
		AdminAPISecret: cfg.GetString(AdminAPISecret),
		Timezone:       cfg.GetString(AppTimezoneKey),
		SeatLockTTL:    cfg.GetDuration(SeatLockTTLKey),
	}
}
