package healthcheckrepo

import (
	domaincache "ticket-reservation/internal/domain/cache"

	"github.com/redis/go-redis/v9"
)

type healthCheckRepositoryImpl struct {
	redisClient redis.UniversalClient
}

func NewHealthCheckRepository(redisClient redis.UniversalClient) domaincache.HealthCheckRepository {
	return &healthCheckRepositoryImpl{
		redisClient: redisClient,
	}
}
