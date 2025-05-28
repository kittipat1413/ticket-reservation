package repository

import (
	"context"
	"ticket-reservation/internal/domain/entity"
	"ticket-reservation/internal/infra/db"

	"github.com/google/uuid"
)

//go:generate mockgen -source=./concert_repository.go -destination=./mocks/concert_repository.go -package=repository_mocks
type ConcertRepository interface {
	CreateOne(ctx context.Context, concert *entity.Concert) (*entity.Concert, error)
	FindOne(ctx context.Context, id uuid.UUID) (*entity.Concert, error)
	WithTx(tx db.SqlExecer) ConcertRepository // Optional: WithTx if you want to use a transaction
}
