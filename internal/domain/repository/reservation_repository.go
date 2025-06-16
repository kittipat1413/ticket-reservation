package repository

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./reservation_repository.go -destination=./mocks/reservation_repository.go -package=repository_mocks
type ReservationRepository interface {
	CreateOne(ctx context.Context, reservation *entity.Reservation) (*entity.Reservation, error)
	FindAll(ctx context.Context, filter FindAllReservationsFilter) (*entity.Reservations, int64, error)
	UpdateOne(ctx context.Context, input UpdateReservationInput) (*entity.Reservation, error)
	WithTx(tx db.SqlExecer) ReservationRepository // Optional: WithTx if you want to use a transaction
}

type FindAllReservationsFilter struct {
	SeatID    *uuid.UUID
	SessionID *string
	Status    *entity.ReservationStatus
	Limit     *int64
	Offset    *int64
}

type UpdateReservationInput struct {
	ID        uuid.UUID
	Status    *entity.ReservationStatus
	ExpiresAt *time.Time
}
