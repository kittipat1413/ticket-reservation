package config

// Service configuration environment variable keys
const (
	ServiceNameKey      = "SERVICE_NAME"
	ServicePortKey      = "SERVICE_PORT"
	ServiceEnvKey       = "SERVICE_ENV"
	ServiceErrPrefixKey = "SERVICE_ERR_PREFIX"
	OtelExporterKey     = "OTEL_EXPORTER_OTLP_ENDPOINT"
	AdminApiKey         = "ADMIN_API_KEY"
	AdminApiSecret      = "ADMIN_API_SECRET"
	AppTimezoneKey      = "APP_TIMEZONE"
)

// Database configuration environment variable keys
const (
	DatabaseUrlKey             = "DATABASE_URL"
	DatabaseMaxOpenConnsKey    = "DATABASE_MAX_OPEN_CONNS"
	DatabaseMaxIdleConnsKey    = "DATABASE_MAX_IDLE_CONNS"
	DatabaseConnMaxLifetimeKey = "DATABASE_CONN_MAX_LIFETIME"  // duration string like "30m"
	DatabaseConnMaxIdleTimeKey = "DATABASE_CONN_MAX_IDLE_TIME" // duration string like "5m"
)

// Default configuration values
// These values are used if the environment variables are not set
var configDefaults = map[string]any{
	// Service configuration
	ServiceNameKey:      "ticket-reservation-api",
	ServicePortKey:      ":8080",
	ServiceEnvKey:       "development",
	AppTimezoneKey:      "Asia/Bangkok",
	ServiceErrPrefixKey: "TR",
	// Database configuration
	DatabaseUrlKey:             "postgres:///ticket-reservation?sslmode=disable",
	DatabaseMaxOpenConnsKey:    30,
	DatabaseMaxIdleConnsKey:    15,
	DatabaseConnMaxLifetimeKey: "30m",
	DatabaseConnMaxIdleTimeKey: "5m",
}
