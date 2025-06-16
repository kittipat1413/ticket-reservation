package repository

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./seat_repository.go -destination=./mocks/seat_repository.go -package=repository_mocks
type SeatRepository interface {
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Seat, error)
	UpdateOne(ctx context.Context, input UpdateSeatInput) (*entity.Seat, error)
	WithTx(tx db.SqlExecer) SeatRepository // Optional: WithTx if you want to use a transaction
}

type UpdateSeatInput struct {
	ID                uuid.UUID
	Status            *entity.SeatStatus
	LockedBySessionID *string
	LockedUntil       *time.Time
}
