package healthcheckrepo

import (
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

type healthCheckRepositoryImpl struct {
	db db.SqlExecer
}

func NewHealthCheckRepository(db db.SqlExecer) repository.HealthCheckRepository {
	return &healthCheckRepositoryImpl{db: db}
}

// WithTx returns a new repository using the provided transaction.
func (r *healthCheckRepositoryImpl) WithTx(tx db.SqlExecer) repository.HealthCheckRepository {
	return &healthCheckRepositoryImpl{db: tx}
}
