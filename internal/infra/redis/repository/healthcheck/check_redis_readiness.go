package healthcheckrepo

import (
	"context"
	"fmt"
	"time"

	errsFramework "github.com/kittipat1413/go-common/framework/errors"
)

func (r *healthCheckRepositoryImpl) CheckRedisReadiness(ctx context.Context) (ok bool, err error) {
	const errLocation = "[repository healthcheck/check_redis_readiness CheckRedisReadiness]"
	defer errsFramework.WrapErrorWithPrefix(errLocation, &err)

	// Test basic connectivity with PING
	pong, err := r.redisClient.Ping(ctx).Result()
	if err != nil {
		return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to ping Redis", err.Error()))
	}

	if pong != "PONG" {
		return false, errsFramework.NewDatabaseError("unexpected Redis PING response", fmt.Sprintf("expected: PONG, got: %s", pong))
	}

	// Test basic operations: SET, GET, DEL
	testKey := fmt.Sprintf("health_check_%d", time.Now().UnixNano())
	testValue := "health_check_value"

	// Test SET operation
	err = r.redisClient.Set(ctx, testKey, testValue, time.Minute).Err()
	if err != nil {
		return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to set key in Redis", err.Error()))
	}

	// Test GET operation
	retrievedValue, err := r.redisClient.Get(ctx, testKey).Result()
	if err != nil {
		// Cleanup on error
		r.redisClient.Del(ctx, testKey)
		return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to get key from Redis", err.Error()))
	}

	if retrievedValue != testValue {
		// Cleanup on error
		r.redisClient.Del(ctx, testKey)
		return false, errsFramework.NewDatabaseError("unexpected Redis GET response", fmt.Sprintf("expected: %s, got: %s", testValue, retrievedValue))
	}

	// Test DEL operation
	deleted, err := r.redisClient.Del(ctx, testKey).Result()
	if err != nil {
		return false, errsFramework.WrapError(err, errsFramework.NewDatabaseError("failed to delete key from Redis", err.Error()))
	}

	if deleted != 1 {
		return false, errsFramework.NewDatabaseError("unexpected Redis DEL response", fmt.Sprintf("expected: 1, got: %d", deleted))
	}

	return true, nil
}
