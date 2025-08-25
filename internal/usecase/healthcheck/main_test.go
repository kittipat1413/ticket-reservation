package usecase_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cache_mocks "ticket-reservation/internal/domain/cache/mocks"
	repository_mocks "ticket-reservation/internal/domain/repository/mocks"
	healthcheckusecase "ticket-reservation/internal/usecase/healthcheck"

	"github.com/kittipat1413/go-common/framework/retry"
)

type testHelper struct {
	ctrl                *gomock.Controller
	retrier             retry.Retrier
	mockDBRepository    *repository_mocks.MockHealthCheckRepository
	mockRedisRepository *cache_mocks.MockHealthCheckRepository
	healthCheckUsecase  healthcheckusecase.HealthCheckUsecase
}

func initTest(t *testing.T) *testHelper {
	ctrl := gomock.NewController(t)

	// Use real retrier with short retry configuration for tests (similar to dependency.go)
	queryBackoff, _ := retry.NewExponentialBackoffStrategy(10*time.Millisecond, 2.0, 100*time.Millisecond)
	retrier, _ := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		Backoff:     queryBackoff,
	})

	mockDBRepository := repository_mocks.NewMockHealthCheckRepository(ctrl)
	mockRedisRepository := cache_mocks.NewMockHealthCheckRepository(ctrl)

	usecase := healthcheckusecase.NewHealthCheckUsecase(
		retrier,
		mockDBRepository,
		mockRedisRepository,
	)

	return &testHelper{
		ctrl:                ctrl,
		retrier:             retrier,
		mockDBRepository:    mockDBRepository,
		mockRedisRepository: mockRedisRepository,
		healthCheckUsecase:  usecase,
	}
}

func (h *testHelper) Done() {
	h.ctrl.Finish()
}

func TestNewHealthCheckUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queryBackoff, _ := retry.NewExponentialBackoffStrategy(10*time.Millisecond, 2.0, 100*time.Millisecond)
	retrier, _ := retry.NewRetrier(retry.Config{
		MaxAttempts: 3,
		Backoff:     queryBackoff,
	})
	mockDBRepo := repository_mocks.NewMockHealthCheckRepository(ctrl)
	mockRedisRepo := cache_mocks.NewMockHealthCheckRepository(ctrl)

	// Execute
	usecase := healthcheckusecase.NewHealthCheckUsecase(retrier, mockDBRepo, mockRedisRepo)

	// Assert
	assert.NotNil(t, usecase)
}
