package reservationrepo

import (
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

type reservationRepositoryImpl struct {
	execer db.SqlExecer
}

func NewReservationRepository(execer db.SqlExecer) repository.ReservationRepository {
	return &reservationRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *reservationRepositoryImpl) WithTx(tx db.SqlExecer) repository.ReservationRepository {
	return &reservationRepositoryImpl{execer: tx}
}
