package zonerepo

import (
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

type zoneRepositoryImpl struct {
	execer db.SqlExecer
}

func NewZoneRepository(execer db.SqlExecer) repository.ZoneRepository {
	return &zoneRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *zoneRepositoryImpl) WithTx(tx db.SqlExecer) repository.ZoneRepository {
	return &zoneRepositoryImpl{execer: tx}
}
