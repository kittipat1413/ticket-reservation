package usecase

import (
	"context"
	"ticket-reservation/internal/domain/repository"

	"github.com/kittipat1413/go-common/framework/retry"
)

//go:generate mockgen -source=./main.go -destination=./mocks/health_check_usecase.go -package=healthcheck_usecasemocks
type HealthCheckUsecase interface {
	CheckReadiness(ctx context.Context) (ok bool, err error)
}

type healthCheckUsecase struct {
	retrier               retry.Retrier
	healthcheckRepository repository.HealthCheckRepository
}

func NewHealthCheckUsecase(
	retrier retry.Retrier,
	healthcheckRepository repository.HealthCheckRepository,
) HealthCheckUsecase {
	return &healthCheckUsecase{
		retrier:               retrier,
		healthcheckRepository: healthcheckRepository,
	}
}
