package seatrepo

import (
	"ticket-reservation/internal/domain/repository"
	"ticket-reservation/internal/infra/db"
)

type seatRepositoryImpl struct {
	execer db.SqlExecer
}

func NewSeatRepository(execer db.SqlExecer) repository.SeatRepository {
	return &seatRepositoryImpl{
		execer: execer,
	}
}

// WithTx returns a new repository with the given transaction
func (r *seatRepositoryImpl) WithTx(tx db.SqlExecer) repository.SeatRepository {
	return &seatRepositoryImpl{
		execer: tx,
	}
}
