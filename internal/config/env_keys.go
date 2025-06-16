package config

// Service configuration environment variable keys
const (
	ServiceNameKey      = "SERVICE_NAME"
	ServicePortKey      = "SERVICE_PORT"
	ServiceEnvKey       = "SERVICE_ENV"
	ServiceErrPrefixKey = "SERVICE_ERR_PREFIX"
	OtelExporterKey     = "OTEL_COLLECTOR_ENDPOINT"
	AdminApiKey         = "ADMIN_API_KEY"    // #nosec G101
	AdminApiSecret      = "ADMIN_API_SECRET" // #nosec G101
	AppTimezoneKey      = "APP_TIMEZONE"
	SeatLockTTLKey      = "SEAT_LOCK_TTL"
)

// Database configuration environment variable keys
const (
	DatabaseUrlKey             = "DATABASE_URL"
	DatabaseMaxOpenConnsKey    = "DATABASE_MAX_OPEN_CONNS"
	DatabaseMaxIdleConnsKey    = "DATABASE_MAX_IDLE_CONNS"
	DatabaseConnMaxLifetimeKey = "DATABASE_CONN_MAX_LIFETIME"  // duration string like "30m"
	DatabaseConnMaxIdleTimeKey = "DATABASE_CONN_MAX_IDLE_TIME" // duration string like "5m"
	RedisAddrsKey              = "REDIS_ADDRS"                 // comma-separated list of Redis addresses
	RedisUsernameKey           = "REDIS_USERNAME"
	RedisPasswordKey           = "REDIS_PASSWORD"      // #nosec G101
	RedisDBKey                 = "REDIS_DB"            // Redis database index
	RedisDialTimeoutKey        = "REDIS_DIAL_TIMEOUT"  // duration string like "5s"
	RedisReadTimeoutKey        = "REDIS_READ_TIMEOUT"  // duration string like "5s"
	RedisWriteTimeoutKey       = "REDIS_WRITE_TIMEOUT" // duration string like "5s"
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
	SeatLockTTLKey:      "300s",
	// Database configuration
	DatabaseUrlKey:             "postgres://postgres:mypass@localhost:5432/ticket-reservation?sslmode=disable",
	DatabaseMaxOpenConnsKey:    30,
	DatabaseMaxIdleConnsKey:    15,
	DatabaseConnMaxLifetimeKey: "30m",
	DatabaseConnMaxIdleTimeKey: "5m",
	RedisAddrsKey:              "localhost:6379",
	RedisUsernameKey:           "",
	RedisPasswordKey:           "", // #nosec G101
	RedisDBKey:                 0,
	RedisDialTimeoutKey:        "3s",
	RedisReadTimeoutKey:        "500ms",
	RedisWriteTimeoutKey:       "500ms",
}
