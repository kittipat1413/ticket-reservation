package usecase

import (
	"context"
	"ticket-reservation/internal/domain/repository"
)

//go:generate mockgen -source=./main.go -destination=./mocks/health_check_usecase.go -package=healthcheck_usecasemocks
type HealthCheckUsecase interface {
	CheckReadiness(ctx context.Context) (ok bool, err error)
}

type healthCheckUsecase struct {
	healthcheckRepository repository.HealthCheckRepository
}

func NewHealthCheckUsecase(
	healthcheckRepository repository.HealthCheckRepository,
) HealthCheckUsecase {
	return &healthCheckUsecase{
		healthcheckRepository: healthcheckRepository,
	}
}
