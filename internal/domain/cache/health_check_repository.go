package cache

import (
	"context"
)

//go:generate mockgen -source=./health_check_repository.go -destination=./mocks/health_check_repository.go -package=cache_mocks
type HealthCheckRepository interface {
	// CheckRedisReadiness verifies connectivity and basic operations
	CheckRedisReadiness(ctx context.Context) (bool, error)
}
