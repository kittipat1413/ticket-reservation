package concertrepo

import (
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

type concertRepositoryImpl struct {
	execer db.SqlExecer
}

func NewConcertRepository(execer db.SqlExecer) repository.ConcertRepository {
	return &concertRepositoryImpl{execer: execer}
}

// WithTx returns a new repository using the provided transaction.
func (r *concertRepositoryImpl) WithTx(tx db.SqlExecer) repository.ConcertRepository {
	return &concertRepositoryImpl{execer: tx}
}
