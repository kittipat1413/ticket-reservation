package repository

import (
	"context"
	"ticket-reservation/internal/infra/db"
)

//go:generate mockgen -source=./health_check_repository.go -destination=./mocks/health_check_repository.go -package=repository_mocks
type HealthCheckRepository interface {
	CheckDatabaseReadiness(ctx context.Context) (ok bool, err error)
	WithTx(tx db.SqlExecer) HealthCheckRepository // Optional: WithTx if you want to use a transaction
}
